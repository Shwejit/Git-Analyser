package main

import "net/http"

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 🔥 allow ALL origins (needed for Chrome extension)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 🔥 allow ALL headers (this is the missing piece)
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// 🔥 allow all common methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 🔥 VERY IMPORTANT: respond to preflight BEFORE hitting routes
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
