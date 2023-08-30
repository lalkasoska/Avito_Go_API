package get_segments

import (
	"avito_go_api/cmd/internal/lib/api/response"
	"avito_go_api/cmd/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	UserId int64 `json:"userId"`
}

type Response struct {
	response.Response
	Segments []string `json:"segments"`
}

type SegmentGetter interface {
	GetSegments(userId int64) ([]string, error)
}

func New(log *slog.Logger, segmentGetter SegmentGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.get_segments.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		segments, err := segmentGetter.GetSegments(req.UserId)
		if err != nil {
			log.Error("failed to get segments", sl.Err(err))
			render.JSON(w, r, response.Error("failed to get segments"))
			return
		}
		log.Info("segments retrieved successfully")
		render.JSON(w, r, Response{response.OK(), segments})
	}
}
