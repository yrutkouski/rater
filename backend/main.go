package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := getEnv("PORT", "8080")
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", allowMethod("GET", securityHeaders(healthHandler)))
	mux.HandleFunc("/api/data", allowMethod("GET", securityHeaders(dataHandler)))

	if staticDir := findStaticDir(); staticDir != "" {
		mux.Handle("/", securityHeadersHandler(spaHandler(staticDir)))
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func securityHeaders(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		handler(w, r)
	}
}

func securityHeadersHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		handler.ServeHTTP(w, r)
	})
}

func allowMethod(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

func findStaticDir() string {
	if dir := os.Getenv("STATIC_DIR"); dir != "" {
		return dir
	}
	for _, dir := range []string{"../frontend/dist", "./frontend/dist"} {
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}
	return ""
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("DB_HOST") == "" {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

type spaHandler string

func (dir spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(string(dir), r.URL.Path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(string(dir), "index.html"))
		return
	}
	http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
}

