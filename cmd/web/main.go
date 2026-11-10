package main

import (
	"database/sql"
	"html/template"
	"log"
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
	errorLog       *log.Logger
	infoLog        *log.Logger
	templateCache  map[string]*template.Template
	users          *models.UserModel
	sessionManager *scs.SessionManager
	formDecoder    *form.Decoder
}

func main() {
	cfg := config.Load()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// db, err := sql.Open("sqlite", "data/app.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	// if err != nil {
	// 	errorLog.Fatal("can't connect")
	// }
	// defer db.Close()

	// if err := db.Ping(); err != nil {
	// 	errorLog.Fatal("can't connect")
	// }

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
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = false // set this to true for https

	formDecoder := form.NewDecoder()
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		// users:         &models.UserModel{DB: db},
		templateCache: templateCache,
	}

	// server := NewServer(app.templateCache, s)
	server := &http.Server{
		Addr:     cfg.Port,
		ErrorLog: errorLog,
		Handler:  app.routes(),
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
