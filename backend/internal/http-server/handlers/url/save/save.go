package save

import (
	"net/http"

	"log/slog"

	resp "project/internal/lib/api/response"
	"project/internal/lib/logger/sl"

	"project/internal/lib/random"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	AliasExists(alias string) (bool, error)
}


type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			log.Info("Generating new alias for URL", slog.String("url", req.URL))
			alias = random.NewRandomString(req.URL, urlSaver, log)
		} else {
			exists, err := urlSaver.AliasExists(alias)
			if err != nil {
				log.Error("Error checking alias existence",
					slog.String("alias", alias),
					sl.Err(err))
				render.JSON(w, r, resp.Error("Error checking alias existence"))
				return
			}

			if exists {
				log.Error("Alias already exists", slog.String("alias", alias))
				render.JSON(w, r, resp.Error("Alias already exists"))
				return
			}
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(alias),
		Alias:    alias,
	})
}
