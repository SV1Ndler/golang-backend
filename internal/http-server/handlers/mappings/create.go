package mappings

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
)

// URLGetter is an interface for getting url by alias.
//

type createMappingRequest struct {
	// File        `json:"image_id" validate:"required"`
	// Created time.Time `json:"created,omitempty"`
}

type createMappingResponse struct {
	resp.Response
	ID int `json:"id,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.40.1 --name=MappingCreater
type MappingCreater interface {
	CreateMapping(imageID int, postID int) (int, error)
}

func Create(log *slog.Logger, mappingCreater MappingCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.images.Create"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		postID, err := strconv.Atoi(chi.URLParam(r, "post-id"))
		if err != nil || postID < 0 {
			// TODO
			log.Info("bad id")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		imageID, err := strconv.Atoi(chi.URLParam(r, "image-id"))
		if err != nil || imageID < 0 {
			// TODO
			log.Info("bad id")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		id, err := mappingCreater.CreateMapping(imageID, postID)
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

		render.JSON(w, r, createMappingResponse{
			Response: resp.OK(),
			ID:       id,
		})
	}
}
