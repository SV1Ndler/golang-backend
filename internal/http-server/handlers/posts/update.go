package posts

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
)

// URLGetter is an interface for getting url by alias.
//

type updatePostRequest struct {
	// ID      int    `json:"id" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	// Created time.Time `json:"created,omitempty"`
}

type UpdatePostResponse struct {
	resp.Response
	ID int `json:"id,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.40.2 --name=PostUpdater
type PostUpdater interface {
	UpdatePost(id int, title string, content string) (int, error)
}

func Update(log *slog.Logger, postUpdater PostUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.posts.Update"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 0 {
			// TODO
			log.Info("bad id")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		var req updatePostRequest

		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty")

			// render.Status(r, http.StatusBadRequest) // ???
			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			// render.Status(r, http.StatusBadRequest) // ???
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		id, err = postUpdater.UpdatePost(id, req.Title, req.Content)
		if err != nil {
			log.Error("failed to update post", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("post updated", slog.Int64("id", int64(id)))

		render.JSON(w, r, UpdatePostResponse{
			Response: resp.OK(),
			ID:       id,
		})
	}
}
