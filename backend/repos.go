package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Repo struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
}

func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	return token
}

func getUserRepos(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token == "" {
		http.Error(w, "Unauthorized", 401)
		return
	}

	req, _ := http.NewRequest(
		"GET",
		"https://api.github.com/user/repos",
		nil,
	)

	// 🔴 THIS IS THE IMPORTANT CHANGE
	req.Header.Set("Authorization", "token "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "GitHub error", 500)
		return
	}
	defer resp.Body.Close()

	var repos []Repo
	json.NewDecoder(resp.Body).Decode(&repos)

	json.NewEncoder(w).Encode(repos)
}
