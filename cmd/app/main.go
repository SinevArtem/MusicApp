package main

import (
	"LoveMusic/internal/config"
	pgAuth "LoveMusic/internal/database/pgsql/auth"
	pg "LoveMusic/internal/database/pgsql/factory"
	pgHomepage "LoveMusic/internal/database/pgsql/page"
	pgTrakcs "LoveMusic/internal/database/pgsql/tracks"
	"LoveMusic/internal/database/redis"
	"LoveMusic/internal/handlers/auth"
	"LoveMusic/internal/handlers/page"
	"LoveMusic/internal/handlers/tracks"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// go run cmd/app/main.go --config=./config/local.yaml
func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	storage, err := pg.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", "error", err)
		os.Exit(1)
	}

	authRepo := pgAuth.NewAuthRepository(storage)
	homepageRepo := pgHomepage.NewHomepageRepository(storage)
	tracksRepo := pgTrakcs.NewTracksRepository(storage)

	redisClient, err := redis.New(&cfg.Redis)
	if err != nil {
		log.Error("failed to init redis", "error", err)
		os.Exit(1)
	}

	fileserver := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.Get("/profile", page.New(log, homepageRepo, redisClient))

	router.Post("/register", auth.RegisterHandler(log, authRepo))
	router.Get("/register", auth.RegisterHandler(log, authRepo))
	router.Get("/login", auth.LoginHandler(log, authRepo, redisClient))
	router.Post("/login", auth.LoginHandler(log, authRepo, redisClient))
	router.Get("/friends", page.UserFriends(log, homepageRepo, redisClient))
	router.Post("/friends", page.UserFriends(log, homepageRepo, redisClient))
	router.Get("/logout", auth.LogoutHandler(log, redisClient))
	router.Get("/collection", page.CollectionHandler(log, homepageRepo, redisClient))
	router.Post("/search_track", tracks.SearchTrack(log, tracksRepo, redisClient))
	router.Get("/search_track", tracks.SearchTrack(log, tracksRepo, redisClient))
	router.Post("/add_track", tracks.AddTrack(log, tracksRepo, redisClient))
	router.Get("/add_track", tracks.AddTrack(log, tracksRepo, redisClient))

	router.Get("/user/{id}", page.UserProfileHandler(log, homepageRepo, redisClient))

	//http.HandleFunc("/profile", LoadProfile)

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", "error", err)
			done <- syscall.SIGTERM
		}

	}()
	log.Info("started server", slog.String("address", cfg.HTTPServer.Address))

	<-done
	log.Info("server shutting down...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown error", "error", err)
	}

	redisClient.Close()
	storage.Close()

	log.Info("server stopped")
}

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
