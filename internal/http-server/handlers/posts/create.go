package posts

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
)

// URLGetter is an interface for getting url by alias.
//

type createPostRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	// Created time.Time `json:"created,omitempty"`
}

type CreatePostResponse struct {
	resp.Response
	ID int `json:"id,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.40.2 --name=PostCreater
type PostCreater interface {
	CreatePost(title string, content string, created time.Time) (int, error)
}

func Create(log *slog.Logger, postCreater PostCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.posts.Create"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req createPostRequest

		err := render.DecodeJSON(r.Body, &req)
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

		id, err := postCreater.CreatePost(req.Title, req.Content, time.Now())
		if err != nil {
			log.Error("failed to create post", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("post created", slog.Int64("id", int64(id)))

		render.JSON(w, r, CreatePostResponse{
			Response: resp.OK(),
			ID:       id,
		})
	}
}
