package get_user_history

import (
	"avito_go_api/cmd/internal/config"
	"avito_go_api/cmd/internal/lib/api/response"
	"avito_go_api/cmd/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"time"
)

type Request struct {
	UserId int64      `json:"userId"`
	Year   int        `json:"year"`
	Month  time.Month `json:"month"`
}

type Response struct {
	response.Response
	Link string `json:"link"`
}

type HistoryGetter interface {
	GetUserHistory(userId int64, year int, month time.Month) error
}

func New(log *slog.Logger, historyGetter HistoryGetter, cfg config.Config) http.HandlerFunc {
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

		err = historyGetter.GetUserHistory(req.UserId, req.Year, req.Month)
		if err != nil {
			log.Error("failed to get user history", sl.Err(err))
			render.JSON(w, r, response.Error("failed to get user history"))
			return
		}
		log.Info("user history retrieved successfully")
		link := cfg.Address + "/report"
		render.JSON(w, r, Response{response.OK(), link})
	}

}
