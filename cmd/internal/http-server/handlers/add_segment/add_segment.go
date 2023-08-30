package add_segment

import (
	"avito_go_api/cmd/internal/lib/api/response"
	"avito_go_api/cmd/internal/lib/logger/sl"
	"avito_go_api/cmd/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	Name string `json:"name"`
}

type SegmentCreator interface {
	CreateSegment(name string) error
}

func New(log *slog.Logger, segmentCreator SegmentCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.add_segment.New"
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
		if req.Name == "" {
			log.Error("segment name cannot be empty", slog.String("name", req.Name))
			render.JSON(w, r, response.Error("segment name cannot be empty"))
			return
		}
		err = segmentCreator.CreateSegment(req.Name)
		if errors.Is(err, storage.ErrSegmentExists) {
			log.Info("segment already exists", slog.String("name", req.Name))
			render.JSON(w, r, response.Error("segment already exists"))
			return
		}
		if err != nil {
			log.Error("failed to create segment", sl.Err(err))
			render.JSON(w, r, response.Error("failed to create segment"))
			return
		}
		log.Info("segment created successfully")
		render.JSON(w, r, response.OK())
	}
}
