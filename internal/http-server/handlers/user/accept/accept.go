package accept

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

type Manager interface {
	AcceptQuest(userId, questId int64) error
}

func New(log *slog.Logger, m Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.accept.New"

		// setting up logger
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// validating userID
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

		// validating quest id
		questID, err := validate.ValidateID(chi.URLParam(r, "questID"))
		if errors.Is(err, validate.ErrEmptyID) {
			log.Info("empty questID")

			render.JSON(w, r, response.Error("empty questID"))

			return
		}
		if errors.Is(err, validate.ErrInvalidID) {
			log.Info("invalid questID")

			render.JSON(w, r, response.Error("invalid questID"))

			return
		}

		// accepting quest
		err = m.AcceptQuest(userID, questID)
		if errors.Is(err, storage.ErrQuestCompleted) {
			log.Info("quest is already done")

			render.JSON(w, r, response.Error("quest is already done"))

			return
		}
		if errors.Is(err, storage.ErrQuestNotFound) {
			log.Info("quest is not found")

			render.JSON(w, r, response.Error("quest not found"))

			return
		}
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found")

			render.JSON(w, r, response.Error("user not found"))

			return
		}

		log.Info("quest accepted")

		render.JSON(w, r,
			response.OK(),
		)

	}
}
