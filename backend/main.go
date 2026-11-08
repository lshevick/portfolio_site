package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// siteDir is the directory containing the static site files.
var siteDir = "../site"

func main() {
	if _, err := os.Stat(siteDir); os.IsNotExist(err) {
		log.Fatalf("site directory not found: %s", siteDir)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	fileServer := http.FileServer(http.Dir(siteDir))
	mux.Handle("/", fileServer)

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}

	absSite, _ := filepath.Abs(siteDir)
	log.Printf("serving site from %s on :%s", absSite, addr)
	log.Fatal(http.ListenAndServe(":"+addr, mux))
}
