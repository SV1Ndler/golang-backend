package images

import (
	"encoding/base64"
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

type createImageRequest struct {
	File    string    `json:"file" validate:"required"`
	Created time.Time `json:"created,omitempty"`
}

type createImageResponse struct {
	resp.Response
	ID int `json:"id,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type ImageCreater interface {
	CreateImage(image []byte, created time.Time) (int, error)
}

func Create(log *slog.Logger, imageCreater ImageCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.images.Create"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req createImageRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty")
			log.Error("\n")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			// log.Error("\n")
			// log.Error(err.Error())
			log.Error("failed to decode request body", sl.Err(err))
			// log.Error("\n")

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

		// alias := req.Alias
		// if alias == "" {
		// 	alias = random.NewRandomString(aliasLength)
		// }

		bytes, _ := base64.StdEncoding.DecodeString(req.File)
		id, err := imageCreater.CreateImage(bytes, req.Created)
		// if errors.Is(err, storage.ErrURLExists) {
		// 	// TODO
		// 	// log.Info("url already exists", slog.String("url", req.URL))

		// 	render.JSON(w, r, resp.Error("url already exists"))

		// 	return
		// }
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", int64(id)))

		render.JSON(w, r, createImageResponse{
			Response: resp.OK(),
			ID:       id,
		})
	}
}
