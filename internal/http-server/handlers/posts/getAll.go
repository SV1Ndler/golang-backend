package posts

import (
	"errors"
	"net/http"
	"time"

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

type getAllPostRequest struct {
}

type getAllPostResponseItem struct {
	ID      int       `json:"id,omitempty"`
	Title   string    `json:"title,omitempty"`
	Content string    `json:"content,omitempty"`
	Created time.Time `json:"created,omitempty"`
}

type getAllPostResponse struct {
	resp.Response
	Array []getAllPostResponseItem `json:"content,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type PostGetterAll interface {
	GetAllPosts() ([]models.Post, error)
}

func GetAll(log *slog.Logger, postGetter PostGetterAll) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.post.Get"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		arr_posts, err := postGetter.GetAllPosts()
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
		resp := getAllPostResponse{Response: resp.OK(), Array: make([]getAllPostResponseItem, 0, len(arr_posts))}
		for idx := range arr_posts {
			resp.Array = append(resp.Array, getAllPostResponseItem{
				ID:      arr_posts[idx].ID,
				Title:   arr_posts[idx].Title,
				Content: arr_posts[idx].Content,
				Created: arr_posts[idx].Created,
			})
		}

		render.JSON(w, r, resp)
	}
}
