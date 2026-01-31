package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func gitPull(repoPath string) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = repoPath
	return cmd.Run()
}

func SyncCommits(repoPath string) error {

	err := gitPull(repoPath)
	if err != nil {
		fmt.Println("Git pull failed:", err)
	} else {
		fmt.Println("Git repository updated")
	}

	cmd := exec.Command("git", "log", "--pretty=format:%H|%ct|%s")
	cmd.Dir = repoPath

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "|", 3)

		if len(parts) == 3 {
			hash := parts[0]
			unixTime := parts[1]
			message := parts[2]

			DB.Exec(
				`INSERT OR IGNORE INTO commits
	(hash, message, commit_time)
	VALUES (?, ?, datetime(?, 'unixepoch'))`,
				hash, message, unixTime,
			)

		}
	}

	return nil
}
