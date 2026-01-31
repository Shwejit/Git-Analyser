package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
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

	cmd = exec.Command(
		"git",
		"log",
		"--name-only",
		"--pretty=format:%ct",
	)
	cmd.Dir = repoPath

	out, _ = cmd.Output()
	scanner = bufio.NewScanner(strings.NewReader(string(out)))

	var currentTime string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		// If line is timestamp
		if _, err := strconv.ParseInt(line, 10, 64); err == nil {
			currentTime = line
			continue
		}

		// Otherwise it's a filename
		DB.Exec(`
        INSERT INTO file_activity
        (file_name, commit_count, last_modified)
        VALUES (?, 1, datetime(?, 'unixepoch'))
        ON CONFLICT(file_name)
        DO UPDATE SET
            commit_count = commit_count + 1,
            last_modified = datetime(?, 'unixepoch')
    `, line, currentTime, currentTime)
	}

	return nil
}
