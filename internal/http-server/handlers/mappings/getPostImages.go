package mappings

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/pkg/models"
)

// URLGetter is an interface for getting url by alias.
//

type getPostImagesRequest struct {
}

type getPostImagesResponseItem struct {
	ID      int       `json:"id,omitempty"`
	Link    string    `json:"link,omitempty"`
	Created time.Time `json:"created,omitempty"`
}

type getPostImagesResponse struct {
	resp.Response
	Array []getPostImagesResponseItem `json:"content,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type ImagesGetterFromPost interface {
	GetPostImages(postID int) ([]models.Image, error)
}

func GetPostImages(log *slog.Logger, imagesGetter ImagesGetterFromPost) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mapping.GetAll"

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

		arr_images, err := imagesGetter.GetPostImages(postID)
		// if errors.Is(err, storage.ErrURLNotFound) {
		// 	//TODO
		// 	// log.Info("url not found", "alias", alias)

		// 	render.JSON(w, r, resp.Error("not found"))

		// 	return
		// }
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}
		// TODO
		// log.Info("got url", slog.String("url", resURL))

		// // redirect to found url
		// http.Redirect(w, r, resURL, http.StatusFound)
		resp := getPostImagesResponse{Response: resp.OK(), Array: make([]getPostImagesResponseItem, 0, len(arr_images))}
		for idx := range arr_images {
			resp.Array = append(resp.Array, getPostImagesResponseItem{
				ID:      arr_images[idx].ID,
				Link:    arr_images[idx].Link,
				Created: arr_images[idx].Created,
			})
		}

		render.JSON(w, r, resp)
	}
}
