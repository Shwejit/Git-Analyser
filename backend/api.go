package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Commit struct {
	ID         int    `json:"id"`
	Hash       string `json:"hash"`
	Message    string `json:"message"`
	CommitTime string `json:"commit_time"`
}

type Stats struct {
	TotalCommits int    `json:"total_commits"`
	LatestCommit string `json:"latest_commit"`
	LastUpdated  string `json:"last_updated"`
}

type ProjectSummary struct {
	TotalFiles    int     `json:"total_files"`
	ActiveFiles   int     `json:"active_files"`
	StableFiles   int     `json:"stable_files"`
	InactiveFiles int     `json:"inactive_files"`
	ActivityScore float64 `json:"activity_score"`
	ProjectState  string  `json:"project_state"`
}

// Health check API
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is running fine 🚀"))
}

// Get commits API
func getCommitsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, hash, message, commit_time FROM commits ORDER BY datetime(commit_time) DESC LIMIT 20")
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var commits []Commit

	for rows.Next() {
		var c Commit
		rows.Scan(&c.ID, &c.Hash, &c.Message, &c.CommitTime)
		commits = append(commits, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commits)
}

func getStatsHandler(w http.ResponseWriter, r *http.Request) {
	var stats Stats

	// 1️⃣ Get total commits
	err := DB.QueryRow("SELECT COUNT(*) FROM commits").
		Scan(&stats.TotalCommits)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// 2️⃣ Get latest commit
	err = DB.QueryRow(
		"SELECT message, commit_time FROM commits ORDER BY datetime(commit_time) DESC LIMIT 1",
	).Scan(&stats.LatestCommit, &stats.LastUpdated)

	if err != nil {
		stats.LatestCommit = "No commits yet"
		stats.LastUpdated = "-"
	}

	// 3️⃣ Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
func getProjectSummary(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`
        SELECT last_modified FROM file_activity
    `)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	total := 0
	active := 0
	stable := 0
	inactive := 0

	for rows.Next() {
		var lastModified string
		rows.Scan(&lastModified)

		t, err := time.Parse(time.RFC3339, lastModified)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04:05", lastModified)
			if err != nil {
				continue
			}
		}

		days := time.Since(t).Hours() / 24

		total++

		if days <= 7 {
			active++
		} else if days <= 30 {
			stable++
		} else {
			inactive++
		}
	}

	activityScore := 0.0
	if total > 0 {
		activityScore = (float64(active) / float64(total)) * 100
	}

	state := "STABLE"
	if activityScore > 50 {
		state = "HIGH ACTIVITY"
	} else if activityScore > 25 {
		state = "EVOLVING"
	}

	summary := ProjectSummary{
		TotalFiles:    total,
		ActiveFiles:   active,
		StableFiles:   stable,
		InactiveFiles: inactive,
		ActivityScore: activityScore,
		ProjectState:  state,
	}

	json.NewEncoder(w).Encode(summary)
}
