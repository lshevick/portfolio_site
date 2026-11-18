package main

import (
	"context"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	if err := run(context.Background(), logger); err != nil {
		logger.Error("fatal error", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg := config.Load()

	s, err := store.New(cfg.DBPath)
	if err != nil {
		return err
	}
	defer s.DB.Close()

	if err := s.MigrateUp(); err != nil {
		return err
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		return err
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

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("starting server", "addr", cfg.Port)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case sig := <-stop:
		logger.Info("shutting down", "signal", sig.String())
	}

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	return server.Shutdown(shutdownCtx)

	// 	slog.Info("starting server", "port", cfg.Port)
	// 	if err := server.ListenAndServe(); err != nil {
	// 		slog.Error("server error", "err", err)
	// 		os.Exit(1)
	// 	}
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
