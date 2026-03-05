package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	args := os.Args[1:]

	if len(args) == 0 {
		if len(cfg.Repos) == 0 {
			printWelcome()
			return
		}
		showCommits(cfg, StartOfToday())
		return
	}

	switch args[0] {
	case "week":
		if len(cfg.Repos) == 0 {
			printWelcome()
			return
		}
		showCommits(cfg, StartOfWeek())

	case "sum":
		if len(cfg.Repos) == 0 {
			printWelcome()
			return
		}
		since := StartOfToday()
		for _, a := range args[1:] {
			if a == "--week" {
				since = StartOfWeek()
				break
			}
		}
		commits := collectCommits(cfg, since)
		PrintCommits(commits)

		if len(commits) == 0 {
			return
		}

		apiKey := GetAPIKey()
		if apiKey == "" {
			fmt.Println()
			fmt.Println("API key not found. Run:")
			fmt.Println(`  clog config --api-key "your-key"`)
			return
		}

		fmt.Println()
		model := GetModel(cfg)
		summary, err := Summarize(commits, apiKey, model)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating summary: %v\n", err)
			os.Exit(1)
		}
		PrintSummaryBlock(summary)

	case "config":
		handleConfig(cfg, args[1:])

	case "repo":
		handleRepo(cfg, args[1:])

	default:
		printUsage()
	}
}

func printWelcome() {
	fmt.Println("Welcome to clog! Let's get you set up.")
	fmt.Println()
	fmt.Println("  Add a repo:       clog repo --add my-project /path/to/repo")
	fmt.Println(`  Set your name:    clog config --author "Your Name"`)
	fmt.Println()
	fmt.Println("  For AI summaries:")
	fmt.Println(`  clog config --api-key "your-key"`)
	fmt.Println("  Then run: clog sum")
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  clog                          Show today's commits")
	fmt.Println("  clog week                     Show this week's commits")
	fmt.Println("  clog sum                       Show today's commits + AI summary")
	fmt.Println("  clog sum --week                Show this week's commits + AI summary")
	fmt.Println(`  clog config --author <name>    Set git author name`)
	fmt.Println(`  clog config --api-key <key>    Set Anthropic API key`)
	fmt.Println(`  clog config --model <model>    Set AI model`)
	fmt.Println("  clog repo --add <name> <path>  Add a repo")
	fmt.Println("  clog repo --remove <name>      Remove a repo")
	fmt.Println("  clog repo --list               List tracked repos")
}

func collectCommits(cfg *Config, since time.Time) []Commit {
	var all []Commit
	for _, repo := range cfg.Repos {
		commits, err := GetCommits(repo.Path, repo.Name, since, cfg.Author)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", repo.Name, err)
			continue
		}
		all = append(all, commits...)
	}
	return all
}

func showCommits(cfg *Config, since time.Time) {
	commits := collectCommits(cfg, since)
	PrintCommits(commits)
}

func handleConfig(cfg *Config, args []string) {
	if len(args) < 2 {
		printConfigUsage()
		return
	}

	switch args[0] {
	case "--author":
		cfg.Author = args[1]
		if err := SaveConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Author set to \"%s\"\n", args[1])

	case "--api-key":
		if err := SaveEnvValue("ANTHROPIC_API_KEY", args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving API key: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("API key saved")

	case "--model":
		cfg.Model = args[1]
		if err := SaveConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Model set to \"%s\"\n", args[1])

	default:
		printConfigUsage()
	}
}

func printConfigUsage() {
	fmt.Println("Usage:")
	fmt.Println(`  clog config --author <name>     Set git author name`)
	fmt.Println(`  clog config --api-key <key>     Set Anthropic API key`)
	fmt.Println(`  clog config --model <model>     Set AI model`)
}

func handleRepo(cfg *Config, args []string) {
	if len(args) == 0 {
		fmt.Println("Usage:")
		fmt.Println("  clog repo --add <name> <path>")
		fmt.Println("  clog repo --remove <name>")
		fmt.Println("  clog repo --list")
		return
	}

	switch args[0] {
	case "--add":
		if len(args) < 3 {
			fmt.Println("Usage: clog repo --add <name> <path>")
			return
		}
		name := args[1]
		path := args[2]
		for _, r := range cfg.Repos {
			if r.Name == name {
				fmt.Printf("Repo %s already exists\n", name)
				return
			}
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("Warning: path does not exist: %s\n", path)
		}
		cfg.Repos = append(cfg.Repos, RepoInfo{Name: name, Path: path})
		if err := SaveConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Added repo \"%s\" at %s\n", name, path)

	case "--remove":
		if len(args) < 2 {
			fmt.Println("Usage: clog repo --remove <name>")
			return
		}
		name := args[1]
		found := false
		var updated []RepoInfo
		for _, r := range cfg.Repos {
			if r.Name == name {
				found = true
				continue
			}
			updated = append(updated, r)
		}
		if !found {
			fmt.Printf("Repo \"%s\" not found\n", name)
			return
		}
		if updated == nil {
			updated = []RepoInfo{}
		}
		cfg.Repos = updated
		if err := SaveConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Removed repo \"%s\"\n", name)

	case "--list":
		if len(cfg.Repos) == 0 {
			fmt.Println("No repos configured. Add one with: clog repo --add <name> /path/to/repo")
			return
		}
		for _, r := range cfg.Repos {
			fmt.Printf("  %s  %s\n", r.Name, r.Path)
		}

	default:
		fmt.Println("Usage:")
		fmt.Println("  clog repo --add <name> <path>")
		fmt.Println("  clog repo --remove <name>")
		fmt.Println("  clog repo --list")
	}
}
