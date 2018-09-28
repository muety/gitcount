package main

import (
	"fmt"
	"sort"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const (
	MAX_SESSION_DIFF_MIN = 120
)

func main() {
	commitMap := make(map[string]*CommitSummary)

	repo, err := git.PlainOpen("/home/ferdinand/dev/mininote")
	CheckError(err)

	// Get all branches
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
	userList := getUsers(commits)
	userMinutes := make(map[string]float64)
	firstCommitAddition := getAverageCommitDiff(commits) * 3

	for _, u := range userList {
		// Use simple heuristic to estimate work
		var minutes float64
		for i := 0; i < len(commits)-1; {
			if commits[i].Email != u {
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
		userMinutes[u] = minutes
		totalMinutes += minutes
	}

	printOut(userMinutes, totalMinutes)
}

func printOut(userMinutes map[string]float64, totalMinutes float64) {
	// TODO: Sort by minutes
	for k, v := range userMinutes {
		fmt.Printf("%s: %.2f hours\n", k, v/60)
	}
	fmt.Printf("---------\nTotal: %.2f hours\n", totalMinutes/60)
}

func getUsers(commits CommitList) []string {
	var userList []string
	userMap := make(map[string]bool)
	for _, c := range commits {
		userMap[c.Email] = true
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
