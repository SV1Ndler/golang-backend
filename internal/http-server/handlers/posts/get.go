package posts

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
	"url-shortener/pkg/models"
)

// URLGetter is an interface for getting url by alias.
//

type getPostRequest struct {
}

type getPostResponse struct {
	resp.Response
	ID      int       `json:"id,omitempty"`
	Title   string    `json:"title,omitempty"`
	Content string    `json:"content,omitempty"`
	Created time.Time `json:"created,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type PostGetter interface {
	GetPost(id int) (models.Post, error)
}

// type PostGetterAll interface {
// 	GetAllPost(title string, content string, created time.Time) (int, error)
// }

func Get(log *slog.Logger, postGetter PostGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.post.Get"

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

		post, err := postGetter.GetPost(id)
		if errors.Is(err, storage.ErrURLNotFound) {
			//TODO
			// log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}
		// TODO
		// log.Info("got url", slog.String("url", resURL))

		// // redirect to found url
		// http.Redirect(w, r, resURL, http.StatusFound)
		render.JSON(w, r, getPostResponse{
			Response: resp.OK(),
			ID:       post.ID,
			Title:    post.Title,
			Content:  post.Content,
			Created:  post.Created,
		})
	}
}
