package images

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

type getImageRequest struct {
}

type getImageResponse struct {
	resp.Response
	ID      int       `json:"id,omitempty"`
	Link    string    `json:"link,omitempty"`
	Created time.Time `json:"created,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type ImageGetter interface {
	GetImage(id int) (models.Image, error)
}

// type PostGetterAll interface {
// 	GetAllPost(title string, content string, created time.Time) (int, error)
// }

func Get(log *slog.Logger, imageGetter ImageGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.images.Get"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// alias := chi.URLParam(r, "alias")
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 0 {
			// TODO
			log.Info("bad id")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		img, err := imageGetter.GetImage(id)
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
		render.JSON(w, r, getImageResponse{
			Response: resp.OK(),
			ID:       img.ID,
			Link:     img.Link,
			Created:  img.Created,
		})
	}
}
