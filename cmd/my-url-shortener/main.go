package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/tuchango/my-url-shortener/internal/config"
	urlDelete "github.com/tuchango/my-url-shortener/internal/http-server/handlers/url/delete"
	"github.com/tuchango/my-url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/tuchango/my-url-shortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/tuchango/my-url-shortener/internal/http-server/middleware/logger"
	"github.com/tuchango/my-url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/tuchango/my-url-shortener/internal/lib/logger/sl"
	"github.com/tuchango/my-url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// init config
	cfg := config.MustLoad()

	// init logger
	log := setupLogger(cfg.Env)

	log.Info("starting my-url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// init storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// id, err := storage.SaveURL("https://google.com", "ggl")
	// if err != nil {
	// 	log.Error("failed to save url", sl.Err(err))
	// 	os.Exit(1)
	// }

	// urlRes, err := storage.GetURL("ggl")
	// if err != nil {
	// 	log.Error("failed to get url", sl.Err(err))
	// 	os.Exit(1)
	// }

	// TODO: init router: chi, chi render
	// TODO: comrare chi and net/http

	router := chi.NewRouter()

	router.Use(middleware.RequestID)

	// TODO: compare 2 loggers
	// router.Use(middleware.Logger)
	// TODO: зачем этот middleware, если всё равно создается log и передается во все функции?
	router.Use(mwLogger.New(log))

	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat) // позволяет оперировать url-паратеметрами в net/http

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("my-url-shortener", map[string]string{
			cfg.HTTPServer.Username: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", urlDelete.New(log, storage)) // FIXME: returns ok when alias not found
	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	// run server

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
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
