package main

import (
	"bufio"
	"os/exec"
	"strings"
	"time"
)

type Commit struct {
	Hash    string
	Message string
	Date    time.Time
	Repo    string
}

func GetCommits(repoPath string, repoName string, since time.Time, author string) ([]Commit, error) {
	sinceStr := since.Format("2006-01-02T15:04:05")
	cmd := exec.Command("git", "-C", repoPath, "log", "--oneline",
		"--since="+sinceStr, "--author="+author, "--format=%H|%s|%ci")

	out, err := cmd.Output()
	if err != nil {
		return []Commit{}, nil
	}

	var commits []Commit
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			continue
		}
		message := parts[1]
		if strings.HasPrefix(message, "Merge") {
			continue
		}
		date, err := time.Parse("2006-01-02 15:04:05 -0700", parts[2])
		if err != nil {
			continue
		}
		commits = append(commits, Commit{
			Hash:    parts[0],
			Message: message,
			Date:    date,
			Repo:    repoName,
		})
	}

	if commits == nil {
		commits = []Commit{}
	}

	return commits, nil
}

func StartOfToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func StartOfWeek() time.Time {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekday := today.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	daysFromMonday := int(weekday) - int(time.Monday)
	return today.AddDate(0, 0, -daysFromMonday)
}
