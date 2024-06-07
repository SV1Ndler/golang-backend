package images

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/pkg/models"
)

// URLGetter is an interface for getting url by alias.
//

type getAllImageRequest struct {
}

type getAllImageResponseItem struct {
	ID      int       `json:"id,omitempty"`
	Link    string    `json:"link,omitempty"`
	Created time.Time `json:"created,omitempty"`
}

type getAllImageResponse struct {
	resp.Response
	Array []getAllImageResponseItem `json:"content,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.40.1 --name=ImageGetterAll
type ImageGetterAll interface {
	GetAllImages() ([]models.Image, error)
}

func GetAll(log *slog.Logger, postGetter ImageGetterAll) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.post.GetAll"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		arr_images, err := postGetter.GetAllImages()
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
		resp := getAllImageResponse{Response: resp.OK(), Array: make([]getAllImageResponseItem, 0, len(arr_images))}
		for idx := range arr_images {
			resp.Array = append(resp.Array, getAllImageResponseItem{
				ID:      arr_images[idx].ID,
				Link:    arr_images[idx].Link,
				Created: arr_images[idx].Created,
			})
		}

		render.JSON(w, r, resp)
	}
}
