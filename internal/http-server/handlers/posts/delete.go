package posts

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

// URLGetter is an interface for getting url by alias.
//

// type deletePostRequest struct {
// }

type deletePostResponse struct {
	resp.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type PostDeleter interface {
	DeletePost(id int) error
}

func Delete(log *slog.Logger, postDeleter PostDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.posts.Delete"

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

		err = postDeleter.DeletePost(id)
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
		render.JSON(w, r, deletePostResponse{
			Response: resp.OK(), 
		})
	}
}
