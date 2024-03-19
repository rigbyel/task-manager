package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rigbyel/task-manager/internal/config"
	createQuest "github.com/rigbyel/task-manager/internal/http-server/handlers/quest/create"
	"github.com/rigbyel/task-manager/internal/http-server/handlers/user/accept"
	createUser "github.com/rigbyel/task-manager/internal/http-server/handlers/user/create"
	"github.com/rigbyel/task-manager/internal/http-server/handlers/user/history"
	"github.com/rigbyel/task-manager/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// initializing config
	cfg := config.MustLoad()

	// initializing logger
	log := setupLogger(cfg.Env)
	log.Info("starting task-manager", slog.String("env", cfg.Env))

	// initializing sqlite storage
	storage, err := storage.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.String("err", err.Error()))
	}

	// intializing chi router
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// handlers
	router.Route("/user", func(r chi.Router) {
		r.Post("/", createUser.New(log, storage.User()))
		r.Get("/{userID}/history", history.New(log, storage.User()))
		r.Post("/{userID}/quests/{questID}", accept.New(log, storage.Manager()))
	})

	router.Post("/quest", createQuest.New(log, storage.Quest()))

	// starting server
	log.Info("starting server", slog.String("addres", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", slog.String("err", err.Error()))
	}

	log.Error("server stopped")

}

// setting up logger
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log

}
