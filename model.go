package main

import (
	"fmt"
	"os"
	"time"
)

type CommitSummary struct {
	Timestamp time.Time
	Email     string
}

type CommitList []*CommitSummary

type MailMapEntry struct {
	ProperEmail string
	CommitEmail string
}

func (l CommitList) Len() int {
	return len(l)
}

func (l CommitList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l CommitList) Less(i, j int) bool {
	return l[i].Timestamp.Before(l[j].Timestamp)
}

func CheckError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

type UserMinute struct {
	Name   string
	Minute float64
}

type UserMinutes []*UserMinute

func (l UserMinutes) Len() int {
	return len(l)
}

func (l UserMinutes) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l UserMinutes) Less(i, j int) bool {
	return l[i].Minute < l[j].Minute
}
