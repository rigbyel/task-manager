package history

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rigbyel/task-manager/internal/lib/api/response"
	"github.com/rigbyel/task-manager/internal/lib/api/validate"
	"github.com/rigbyel/task-manager/internal/storage"
)

type Response struct {
	response.Response
	QuestList []string `json:"quest_list"`
	Balance   int      `json:"balance"`
}

type HistoryGetter interface {
	History(userID int64) ([]string, error)
	Balance(userID int64) (int, error)
}

func New(log *slog.Logger, h HistoryGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.history.New"

		// setting up logger
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// validating user id
		userID, err := validate.ValidateID(chi.URLParam(r, "userID"))
		if errors.Is(err, validate.ErrEmptyID) {
			log.Info("empty userID")

			render.JSON(w, r, response.Error("empty userID"))

			return
		}
		if errors.Is(err, validate.ErrInvalidID) {
			log.Info("invalid userID")

			render.JSON(w, r, response.Error("invalid userID"))

			return
		}

		// getting user history
		questList, err := h.History(userID)
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found", slog.Int64("userID", userID))

			render.JSON(w, r, response.Error("user not found"))

			return
		}
		if err != nil {
			log.Error("error getting user history", slog.String("error", err.Error()))
		}

		log.Info("user history accessed")

		// getting user balance
		balance, err := h.Balance(userID)
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found", slog.Int64("userID", userID))

			render.JSON(w, r, response.Error("user not found"))

			return
		}
		if err != nil {
			log.Error("error getting user balance", slog.String("error", err.Error()))
		}

		log.Info("user balance accessed")

		render.JSON(w, r,
			Response{
				Response:  response.OK(),
				QuestList: questList,
				Balance:   balance,
			},
		)
	}
}
