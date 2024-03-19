package create

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	response "github.com/rigbyel/task-manager/internal/lib/api/response"
	"github.com/rigbyel/task-manager/internal/model"
	"github.com/rigbyel/task-manager/internal/storage"
)

type Request struct {
	Name string `json:"name" validate:"required"`
	Cost int    `json:"cost" validate:"required"`
}

type Response struct {
	response.Response
	ID int64 `json:"id"`
}

type QuestCreator interface {
	CreateQuest(u *model.Quest) (*model.Quest, error)
}

func New(log *slog.Logger, questCreator QuestCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.quest.create.New"

		// setting up logger
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		// decoding request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("failed to decode request body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// validating request
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", slog.String("error", err.Error()))

			render.JSON(w, r, response.ValidationErrors(validateErr))

			return
		}

		// creating quest
		quest := &model.Quest{
			Name: req.Name,
			Cost: req.Cost,
		}

		quest, err = questCreator.CreateQuest(quest)
		if errors.Is(err, storage.ErrQuestExists) {
			log.Info("quest already exists", slog.String("quest", req.Name))

			render.JSON(w, r, response.Error("quest already exists"))

			return
		}
		if err != nil {
			log.Error("error creating quest", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("error creating quest"))

			return
		}

		log.Info("quest added", slog.Int64("id", quest.Id))

		render.JSON(w, r,
			Response{
				Response: response.OK(),
				ID:       quest.Id,
			},
		)
	}
}
