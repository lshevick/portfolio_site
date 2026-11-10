package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	uiDir := os.Getenv("UI_DIR")
	if uiDir == "" {
		uiDir = "ui"
	}

	t, err := loadTemplates(uiDir)
	if err != nil {
		log.Fatalf("loading templates: %v", err)
	}
	templates = t

	handler := buildHandler(uiDir)

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}

	staticDir := filepath.Join(uiDir, "static")
	log.Printf("serving pages from %s on :%s (static assets: %s)", uiDir, addr, staticDir)
	if err := http.ListenAndServe(":"+addr, handler); err != nil {
		log.Fatal(err)
	}
}
