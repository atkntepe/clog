package main

import (
	"fmt"
	"strings"
)

func PrintCommits(commits []Commit) {
	if len(commits) == 0 {
		return
	}

	repoOrder := []string{}
	grouped := map[string][]Commit{}
	for _, c := range commits {
		if _, exists := grouped[c.Repo]; !exists {
			repoOrder = append(repoOrder, c.Repo)
		}
		grouped[c.Repo] = append(grouped[c.Repo], c)
	}

	for i, repo := range repoOrder {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("● %s\n", repo)
		for _, c := range grouped[repo] {
			hash := c.Hash
			if len(hash) > 7 {
				hash = hash[:7]
			}
			fmt.Printf("  %s  %s\n", hash, c.Message)
		}
	}

	fmt.Println()
	fmt.Printf("  %d commits across %d repos\n", len(commits), len(repoOrder))
}

func PrintSummaryBlock(summary string) {
	separator := "─────────────────────────────────────────"
	fmt.Println(separator)
	fmt.Println("  AI Summary")
	fmt.Println()
	for _, line := range strings.Split(summary, "\n") {
		fmt.Printf("  %s\n", line)
	}
	fmt.Println(separator)
}
