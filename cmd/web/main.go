package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lshevick/internal/config"
	"github.com/lshevick/internal/models"
	"github.com/lshevick/internal/store"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "modernc.org/sqlite"
)

type application struct {
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	users          *models.UserModel
	sessionManager *scs.SessionManager
	formDecoder    *form.Decoder
	store          *store.Store
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := config.Load()

	s, err := store.New(cfg.DBPath)
	if err != nil {
		slog.Error("opening store", "err", err)
		os.Exit(1)
	}
	if err := s.Migrate(); err != nil {
		slog.Error("running migrations", "err", err)
		os.Exit(1)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		slog.Error("error loading teplatecache...", "err", err)
		return
	}

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(s.DB)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = false // set this to true for https

	app := &application{
		logger:         logger,
		users:          &models.UserModel{DB: s.DB},
		templateCache:  templateCache,
		sessionManager: sessionManager,
		formDecoder:    form.NewDecoder(),
		store:          s,
	}

	// server := NewServer(app.templateCache, s)
	server := &http.Server{
		Addr:    cfg.Port,
		Handler: app.routes(),
	}

	slog.Info("starting server", "port", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
	//////////
	// uiDir := os.Getenv("UI_DIR")
	// if uiDir == "" {
	// 	uiDir = "ui"
	// }

	// t, err := loadTemplates(uiDir)
	// if err != nil {
	// 	log.Fatalf("loading templates: %v", err)
	// }
	// templates = t

	// handler := buildHandler(uiDir)

	// addr := os.Getenv("PORT")
	// if addr == "" {
	// 	addr = "8080"
	// }

	// staticDir := filepath.Join(uiDir, "static")
	// log.Printf("serving pages from %s on :%s (static assets: %s)", uiDir, addr, staticDir)
	// if err := http.ListenAndServe(":"+addr, handler); err != nil {
	// 	log.Fatal(err)
	// }
}
