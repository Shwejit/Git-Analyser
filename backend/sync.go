package main

import (
	"net/http"
)

func syncHandler(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	repo := r.URL.Query().Get("repo")
	token := extractToken(r) // uses function from repos.go

	if owner == "" || repo == "" || token == "" {
		http.Error(w, "Missing owner / repo / token", http.StatusBadRequest)
		return
	}

	err := SyncFromGitHub(owner, repo, token)
	if err != nil {
		http.Error(w, "Sync failed", http.StatusInternalServerError)
		return
	}

	saveSnapshot(repo)

	w.Write([]byte("Synced successfully"))
}

// ---------- SNAPSHOT ----------
func saveSnapshot(repo string) {
	row := DB.QueryRow(`
		SELECT
			SUM(CASE WHEN julianday('now') - julianday(last_modified) <= 7 THEN 1 ELSE 0 END),
			SUM(CASE WHEN julianday('now') - julianday(last_modified) BETWEEN 7 AND 30 THEN 1 ELSE 0 END),
			SUM(CASE WHEN julianday('now') - julianday(last_modified) > 30 THEN 1 ELSE 0 END)
		FROM file_activity
	`)

	var active, stable, inactive int
	row.Scan(&active, &stable, &inactive)

	total := active + stable + inactive
	score := 0.0
	if total > 0 {
		score = (float64(active) / float64(total)) * 100
	}

	DB.Exec(`
		INSERT INTO repo_snapshots
		(repo_name, active_files, stable_files, inactive_files, activity_score)
		VALUES (?, ?, ?, ?, ?)
	`, repo, active, stable, inactive, score)
}
