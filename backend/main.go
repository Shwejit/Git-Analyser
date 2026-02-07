package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	err := InitDB()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/project/summary", getProjectSummary)
	http.HandleFunc("/sync", syncHandler)
	http.HandleFunc("/repos", getUserRepos)
	http.HandleFunc("/history", getRepoHistory)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Backend running on port", port)

	// ✅ Enable CORS wrapper
	handler := enableCORS(http.DefaultServeMux)

	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		panic(err)
	}
}
