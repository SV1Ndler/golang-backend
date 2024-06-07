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

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
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
			name:      "Delete Error",
			id:        60,
			respError: "internal error",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			postDeleterMock := mocks.NewPostDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				postDeleterMock.On("DeletePost", tc.id).
					Return(tc.mockError).Once()
			}

			handler := posts.Delete(slogdiscard.NewDiscardLogger(), postDeleterMock)

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

			var resp posts.DeletePostResponse

			require.NoError(t, json.Unmarshal(respBody, &resp))
			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
