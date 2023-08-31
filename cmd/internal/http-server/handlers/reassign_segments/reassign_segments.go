package reassign_segments

import (
	"avito_go_api/cmd/internal/lib/api/response"
	"avito_go_api/cmd/internal/lib/logger/sl"
	"avito_go_api/cmd/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strings"
)

type Request struct {
	SegmentsToAdd    []string `json:"segmentsToAdd"`
	SegmentsToRemove []string `json:"segmentsToRemove"`
	UserId           int64    `json:"userId"`
}

type SegmentAssigner interface {
	ReassignSegments(addSegments []string, removeSegments []string, userId int64) error
}

func New(log *slog.Logger, segmentAssigner SegmentAssigner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.reassign_segments.New"
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

		if len(req.SegmentsToRemove) == 0 && len(req.SegmentsToAdd) == 0 {
			log.Error("no segments were given")
			render.JSON(w, r, response.Error("no segments were given"))
			return
		}

		err = segmentAssigner.ReassignSegments(req.SegmentsToAdd, req.SegmentsToRemove, req.UserId)
		if errors.Is(err, storage.ErrUserHasSegment) {
			log.Info("user already has some of the segments", slog.String("segmentsToAdd", strings.Join(req.SegmentsToAdd, ",")))
			render.JSON(w, r, response.Error("user already has some of the segments"))
			return
		}
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found", slog.Int64("userId", req.UserId))
			render.JSON(w, r, response.Error("user not found"))
			return
		}
		if err != nil {
			log.Error("failed to reassign segments", sl.Err(err))
			render.JSON(w, r, response.Error("failed to reassign segments"))
			return
		}
		log.Info("segments reassigned successfully")
		render.JSON(w, r, response.OK())
	}
}
