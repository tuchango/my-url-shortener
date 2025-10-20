package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/tuchango/my-url-shortener/internal/lib/api/response"
	"github.com/tuchango/my-url-shortener/internal/lib/logger/sl"
	"github.com/tuchango/my-url-shortener/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v3@v3.5.5 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request")) // TODO: better msg

			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", "alias", alias)

				render.JSON(w, r, resp.Error("not found")) // TODO: better msg

				return
			}

			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error")) // TODO: как понять, когда не нужно давать инфу на фронт, а когда нужно? (тут не нужно)

			return
		}

		log.Info("got url", slog.String("url", resURL)) // TODO: why 'slog.String("url", resURL)', not '"url", resURL'

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
