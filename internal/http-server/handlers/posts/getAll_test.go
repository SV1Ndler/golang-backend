package posts_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/posts"
	"url-shortener/internal/http-server/handlers/posts/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/pkg/models"
)

func TestGetAllHandler(t *testing.T) {
	cases := []struct {
		name      string
		respError string
		mockError error
	}{
		{
			name: "Success",
		},
		{
			name:      "GetAll Error",
			respError: "internal error",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			postGetterAllMock := mocks.NewPostGetterAll(t)

			if tc.respError == "" || tc.mockError != nil {
				postGetterAllMock.On("GetAllPosts").
					Return([]models.Post{{ID: 1, Title: "title", Content: "c"}}, tc.mockError).Once()
			}

			handler := posts.GetAll(slogdiscard.NewDiscardLogger(), postGetterAllMock)

			req, err := http.NewRequest(http.MethodPost, "/posts", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp posts.GetAllPostResponse

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
