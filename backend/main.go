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

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/project/summary", getProjectSummary)
	mux.HandleFunc("/auth/github", githubLogin)
	mux.HandleFunc("/auth/callback", githubCallback)
	mux.HandleFunc("/sync", syncHandler)
	mux.HandleFunc("/repos", getUserRepos)
	mux.HandleFunc("/history", getRepoHistory)

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("value")
		fmt.Fprintf(w, `
		<script>
		window.opener.postMessage({ token: "%s" }, "*");
		window.close();
		</script>
		`, token)
	})

	// 🔥 GLOBAL HANDLER (THIS FIXES CORS 100%)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow Chrome extension to talk to backend
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 🔥 Handle ALL preflight requests HERE
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass real requests to router
		mux.ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Backend running on port", port)
	panic(http.ListenAndServe(":"+port, handler))
}
