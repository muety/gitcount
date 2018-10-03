package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const (
	MAX_SESSION_DIFF_MIN = 120
)

func main() {
	// Initialization stuff
	cwd, err := os.Getwd()
	CheckError(err)

	dirPtr := flag.String("dir", cwd, "Project root directory absolute path")
	flag.Parse()

	fmt.Printf("Project root: %s\n", *dirPtr)
	repo, err := git.PlainOpen(*dirPtr)
	CheckError(err)

	var mailmap map[string]*MailMapEntry

	if _, err := os.Stat(filepath.Join(*dirPtr, ".mailmap")); !os.IsNotExist(err) {
		mailmapFile, err := os.Open(filepath.Join(*dirPtr, ".mailmap"))
		defer mailmapFile.Close()
		mailmap, err = ReadMailMap(mailmapFile)
		CheckError(err)
	}

	// Get all branches
	commitMap := make(map[string]*CommitSummary)
	branches, err := repo.Branches()
	CheckError(err)

	// Extract all commits fom all branches, uniquely
	err = branches.ForEach(func(b *plumbing.Reference) error {
		head, err := repo.Head()
		CheckError(err)

		commits, err := repo.Log(&git.LogOptions{From: head.Hash()})
		err = commits.ForEach(func(c *object.Commit) error {
			commitMap[c.Hash.String()] = &CommitSummary{
				Timestamp: c.Committer.When,
				Email:     c.Committer.Email,
			}
			return nil
		})
		return nil
	})
	CheckError(err)

	// Sort commits by time
	commits := make(CommitList, 0, len(commitMap))
	for _, v := range commitMap {
		commits = append(commits, v)
	}
	sort.Sort(commits)

	var totalMinutes float64
	var userMinutes UserMinutes
	userList := getUsers(mailmap, commits)
	firstCommitAddition := getAverageCommitDiff(commits) * 3

	for _, u := range userList {
		// Use simple heuristic to estimate work
		var minutes float64
		for i := 0; i < len(commits)-1; {
			if commitEmail := commits[i].Email; getEmail(mailmap, commitEmail) != u {
				i += 1
				continue
			}
			diff := commits[i+1].Timestamp.Sub(commits[i].Timestamp)
			if diff.Minutes() <= MAX_SESSION_DIFF_MIN {
				minutes += diff.Minutes()
			} else {
				minutes += firstCommitAddition
			}
			i += 1
		}
		userMinutes = append(userMinutes, &UserMinute{
			Name:   getEmail(mailmap, u),
			Minute: minutes,
		})
		totalMinutes += minutes
	}

	printOut(userMinutes, totalMinutes)
}

func printOut(userMinutes UserMinutes, totalMinutes float64) {
	sort.Sort(sort.Reverse(userMinutes))
	for _, u := range userMinutes {
		fmt.Printf("%s: %.2f hours\n", u.Name, u.Minute/60)
	}
	fmt.Printf("---------\nTotal: %.2f hours\n", totalMinutes/60)
}

func getUsers(mailmap map[string]*MailMapEntry, commits CommitList) []string {
	var userList []string
	userMap := make(map[string]bool)
	for _, c := range commits {
		userMap[getEmail(mailmap, c.Email)] = true
	}
	for k := range userMap {
		userList = append(userList, k)
	}
	return userList
}

func getAverageCommitDiff(commits CommitList) float64 {
	var minutes float64
	var count int
	for i := 0; i < len(commits)-1; {
		diff := commits[i+1].Timestamp.Sub(commits[i].Timestamp)
		if diff.Minutes() <= MAX_SESSION_DIFF_MIN {
			minutes += diff.Minutes()
			count += 1
		}
		i += 1
	}
	return minutes / float64(count)
}

func getEmail(mailmap map[string]*MailMapEntry, email string) string {
	if entry := mailmap[email]; entry != nil {
		return entry.ProperEmail
	} else {
		return email
	}
}
