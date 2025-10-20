package delete

import (
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/tuchango/my-url-shortener/internal/lib/api/response"
	"github.com/tuchango/my-url-shortener/internal/lib/logger/sl"
	"github.com/tuchango/my-url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v3@v3.5.5 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

// TODO: make tests
func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		// create logger
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		if err := urlDeleter.DeleteURL(alias); err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("alias not found", "alias", alias)

				w.WriteHeader(http.StatusNoContent)

				return
			}

			log.Error("failed to delete url", sl.Err(err))

			w.WriteHeader(http.StatusNoContent)

			return
		}

		log.Info("got url")

		w.WriteHeader(http.StatusNoContent)

	}
}
