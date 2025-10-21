package save

import (
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/tuchango/my-url-shortener/internal/lib/api/response"
	"github.com/tuchango/my-url-shortener/internal/lib/logger/sl"
	"github.com/tuchango/my-url-shortener/internal/lib/random"
	"github.com/tuchango/my-url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// TODO: move to config
const aliasLength = 6

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: разобраться, зачем нужен интерфейс, и почему он именно здесь
//
//go:generate go run github.com/vektra/mockery/v3@v3.5.5 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		// create logger
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// decode request json
		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// validate request
		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))

			validationErrors := err.(validator.ValidationErrors)
			render.JSON(w, r, resp.ValidationError(validationErrors))

			return
		}

		alias := req.Alias
		if alias == "" {
			// TODO: обход/проверка коллизий. мб есть пост в тгк Тузова
			alias = random.GenerateRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				// TODO: тут есть url? а должен быть?
				log.Info("url already exists", sl.Err(err))

				// TODO: а фронт знает к какому url? или ему не нужно знать?
				render.JSON(w, r, resp.Error("url already exists"))

				return
			}

			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
