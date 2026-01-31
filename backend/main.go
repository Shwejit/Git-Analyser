package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	err := InitDB()
	if err != nil {
		fmt.Println("DB error:", err)
		return
	}
	fmt.Println("Database ready")

	repoPath := "D:/LeetCode"

	err = SyncCommits(repoPath)
	if err != nil {
		fmt.Println("Sync error:", err)
		return
	}
	fmt.Println("Initial sync done")

	// Start API server
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/commits", getCommitsHandler)
	http.HandleFunc("/stats", getStatsHandler)

	go func() {
		println("🚀 API running on http://localhost:8080")
		http.ListenAndServe(":8080", nil)
	}()

	for {
		time.Sleep(30 * time.Second)

		fmt.Println("Syncing Git changes...", time.Now())
		err := SyncCommits(repoPath)
		if err != nil {
			fmt.Println("Sync error:", err)
		} else {
			fmt.Println("Sync complete")
		}
	}
}
