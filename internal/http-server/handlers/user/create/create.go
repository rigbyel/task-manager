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
	Name    string `json:"name" validate:"required"`
	Balance int    `json:"balance,omitempty"`
}

type Response struct {
	response.Response
	ID int64 `json:"id"`
}

type UserCreator interface {
	CreateUser(u *model.User) (*model.User, error)
}

func New(log *slog.Logger, userCreator UserCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.create.New"

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

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", slog.String("error", err.Error()))

			render.JSON(w, r, response.ValidationErrors(validateErr))

			return
		}

		user := &model.User{
			Name:    req.Name,
			Balance: req.Balance,
		}

		// creating user
		user, err = userCreator.CreateUser(user)
		if errors.Is(err, storage.ErrUserExists) {
			log.Info("user already exists", slog.String("user", req.Name))

			render.JSON(w, r, response.Error("user already exists"))

			return
		}
		if err != nil {
			log.Error("error creating user", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("error creating user"))

			return
		}

		log.Info("user added", slog.Int64("id", user.Id))

		render.JSON(w, r,
			Response{
				Response: response.OK(),
				ID:       user.Id,
			},
		)
	}
}
