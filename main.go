package main

import (
	"fmt"
	"sort"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const (
	MAX_SESSION_DIFF_MIN      = 120
	FIRST_COMMIT_ADDITION_MIN = 120
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

	// Use simple heuristic to estimate work
	var totalMinutes float64
	for i := 0; i < len(commits)-1; {
		diff := commits[i+1].Timestamp.Sub(commits[i].Timestamp)
		if diff.Minutes() <= MAX_SESSION_DIFF_MIN {
			totalMinutes += diff.Minutes()
		} else {
			totalMinutes += FIRST_COMMIT_ADDITION_MIN
		}
		i += 1
	}

	fmt.Println(totalMinutes)
}
