package mappings

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/pkg/models"
)

// URLGetter is an interface for getting url by alias.
//

type getAllMappingRequest struct {
}

type getAllMappingResponseItem struct {
	ImageID int    `json:"image-id,omitempty"`
	Link    string `json:"link,omitempty"`
}

type getAllMappingResponse struct {
	resp.Response
	Array []getAllMappingResponseItem `json:"content,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type MappingGetterAll interface {
	GetAllMappingsWithLink() ([]models.PostImageMappingWithLink, error)
}

func GetAll(log *slog.Logger, mappingGetter MappingGetterAll) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mapping.GetAll"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		arr_mappings, err := mappingGetter.GetAllMappingsWithLink()
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
		resp := getAllMappingResponse{Response: resp.OK(), Array: make([]getAllMappingResponseItem, 0, len(arr_mappings))}
		for idx := range arr_mappings {
			resp.Array = append(resp.Array, getAllMappingResponseItem{
				ImageID: arr_mappings[idx].ImageID,
				Link:    arr_mappings[idx].Link,
			})
		}

		render.JSON(w, r, resp)
	}
}
