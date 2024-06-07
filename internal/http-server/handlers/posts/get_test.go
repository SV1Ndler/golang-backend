package posts_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"url-shortener/internal/http-server/handlers/posts"
	"url-shortener/internal/http-server/handlers/posts/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/pkg/models"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {
	cases := []struct {
		name      string
		id        int
		respError string
		mockError error
	}{
		{
			name: "Success",
			id:   6,
		},
		{
			name:      "Get Error",
			id:        60,
			respError: "internal error",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			postGetterMock := mocks.NewPostGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				postGetterMock.On("GetPost", tc.id).
					Return(models.Post{ID: 1, Title: "title", Content: "c"}, tc.mockError).Once()
			}

			handler := posts.Get(slogdiscard.NewDiscardLogger(), postGetterMock)

			r := chi.NewRouter()
			r.Delete("/posts/{id}", handler)

			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodDelete, ts.URL+"/posts/"+strconv.Itoa(tc.id), nil)
			require.NoError(t, err)

			respRaw, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, respRaw.StatusCode, http.StatusOK)

			respBody, err := io.ReadAll(respRaw.Body)
			require.NoError(t, err)

			var resp posts.GetPostResponse

			require.NoError(t, json.Unmarshal(respBody, &resp))
			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
