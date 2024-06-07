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
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/exp/slog"

	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/images"
	"url-shortener/internal/http-server/handlers/login"
	"url-shortener/internal/http-server/handlers/mappings"
	"url-shortener/internal/http-server/handlers/posts"

	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/pkg/models/sql"
	// "url-shortener/pkg/models/simple"
	// "github.com/go-chi/cors"
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

	// clientID := "01e493a6685580c" TODO
	storage, er := sql.New()
	// storagePosts := simple.NewPost()
	// storageImages, er := simple.NewImage(clientID)
	// storageMapping := simple.NewPostImageMapping(storageImages)
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

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "UPDATE", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Group(func(r chi.Router) {
		// var tokenAuth *jwtauth.JWTAuth
		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/posts", posts.Create(log, storage))
		r.Delete("/posts/{id}", posts.Delete(log, storage))
		r.Patch("/posts/{id}", posts.Update(log, storage))

		r.Delete("/images/{id}", images.Delete(log, storage))
		r.Post("/images", images.Create(log, storage))

		r.Post("/posts/{post-id}/images/{image-id}", mappings.Create(log, storage))
		r.Get("/posts/{post-id}/images", mappings.GetPostImages(log, storage))
	})

	router.Get("/posts", posts.GetAll(log, storage))
	router.Get("/posts/{id}", posts.Get(log, storage))
	// router.Delete("/posts/{id}", posts.Delete(log, storage))
	// router.Post("/posts", posts.Create(log, storage))
	// router.Patch("/posts/{id}", posts.Update(log, storage))

	router.Get("/images/{id}", images.Get(log, storage))
	router.Get("/images", images.GetAll(log, storage))
	// router.Delete("/images/{id}", images.Delete(log, storage))
	// router.Post("/images", images.Create(log, storage))

	// router.Post("/posts/{post-id}/images/{image-id}", mappings.Create(log, storage))
	// router.Get("/posts/{post-id}/images", mappings.GetPostImages(log, storage))

	router.Get("/login", login.New(log, storage))

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
