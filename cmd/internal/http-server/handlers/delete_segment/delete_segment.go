package delete_segment

import (
	"avito_go_api/cmd/internal/lib/api/response"
	"avito_go_api/cmd/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	Name string `json:"name"`
}

type SegmentDeleter interface {
	DeleteSegment(name string) error
}

func New(log *slog.Logger, segmentDeleter SegmentDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete_segment.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		err = segmentDeleter.DeleteSegment(req.Name)
		if err != nil {
			log.Error("failed to delete segment", sl.Err(err))
			render.JSON(w, r, response.Error("failed to delete segment"))
			return
		}
		log.Info("segment deleted successfully")
		render.JSON(w, r, response.OK())
	}
}
