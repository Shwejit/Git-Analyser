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

	// 🔥 CREATE OUR OWN ROUTER (NO DEFAULT MUX)
	mux := http.NewServeMux()

	// Routes
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Backend running on port", port)

	// 🔥 THIS LINE IS THE REAL FIX
	handler := enableCORS(mux)

	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		panic(err)
	}
}
