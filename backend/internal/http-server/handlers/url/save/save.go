package save


import (
	"net/http"

	"log/slog"

	resp "project/internal/lib/api/response"
	"project/internal/lib/logger/sl"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
)

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}


type Request struct {
	URL string `json:"url" validate:"required,url"`
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
			slog.String("request_id", r.Context().Value(middleware.RequestIDKey).(string)),
		)

		var req Request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request", sl.Err(err))
		}

	}
}