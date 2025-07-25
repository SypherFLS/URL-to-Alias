package out

import (
	"net/http"

	"log/slog"

	resp "project/internal/lib/api/response"
	"project/internal/lib/logger/sl"
	"project/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetAllUrls() ([]storage.URLRecord, error)
}

type PieceOfOutput struct {
	URL      string `json:"url"`
	ShortURL string `json:"short_url"`
}

type Response struct {
	resp.Response
	Data []PieceOfOutput `json:"data"`
}

func TakeAllUrls(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.out.TakeAllUrls"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("attempting to get all URLs")

		records, err := urlGetter.GetAllUrls()
		if err != nil {
			log.Error("failed to get URLs", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get URLs"))
			return
		}

		log.Info("successfully retrieved URLs", slog.Int("count", len(records)))


		var output []PieceOfOutput
		for _, record := range records {
			output = append(output, PieceOfOutput{
				URL:      record.URL,
				ShortURL: record.Alias,
			})
		}

		log.Info("URLs converted to output format", slog.Int("count", len(output)))

		responseOK(w, r, output)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, data []PieceOfOutput) {
	render.JSON(w, r, Response{
		Response: resp.OK("URLs retrieved successfully"),
		Data:     data,
	})
}
