package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"

	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/images"
	"url-shortener/internal/http-server/handlers/mappings"
	"url-shortener/internal/http-server/handlers/posts"

	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/pkg/models/simple"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting url-shortener",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	clientID := "01e493a6685580c"
	storagePosts := simple.NewPost()
	storageImages, er := simple.NewImage(clientID)
	storageMapping := simple.NewPostImageMapping(storageImages)
	if er != nil {
		log.Error("failed to init storage", sl.Err(er))
		os.Exit(1)
	}
	// storage, err := sqlite.New(cfg.StoragePath)
	// if err != nil {
	// 	log.Error("failed to init storage", sl.Err(err))
	// 	os.Exit(1)
	// }

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// router.Route("/url", func(r chi.Router) {
	// 	r.Use(middleware.BasicAuth("url-shortener", map[string]string{
	// 		cfg.HTTPServer.User: cfg.HTTPServer.Password,
	// 	}))

	// 	r.Post("/", save.New(log, storage))
	// 	// TODO: add DELETE /url/{id}
	// })

	// router.Get("/{alias}", redirect.New(log, storage))
	router.Get("/posts/{id}", posts.Get(log, storagePosts))
	router.Delete("/posts/{id}", posts.Delete(log, storagePosts))
	router.Get("/posts", posts.GetAll(log, storagePosts))
	router.Post("/posts", posts.Create(log, storagePosts))

	router.Get("/images/{id}", images.Get(log, storageImages))
	router.Delete("/images/{id}", images.Delete(log, storageImages))
	router.Get("/images", images.GetAll(log, storageImages))
	router.Post("/images", images.Create(log, storageImages))

	router.Get("/posts/{post-id}/images/{image-id}", mappings.Create(log, storageMapping))
	router.Get("/posts/{post-id}/images", mappings.GetAll(log, storageMapping))

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	// TODO: close storage

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
