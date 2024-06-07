package posts_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/posts"
	"url-shortener/internal/http-server/handlers/posts/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
)

func TestCreateHandler(t *testing.T) {
	cases := []struct {
		name      string
		title     string
		content   string
		respError string
		mockError error
	}{
		{
			name:    "Success",
			title:   "Title",
			content: "Many samsing words.",
		},
		{
			name:      "Empty title",
			title:     "",
			content:   "Many samsing words.",
			respError: "field Title is a required field",
		},
		{
			name:      "Empty content",
			title:     "Title",
			content:   "",
			respError: "field Content is a required field",
		},
		{
			name:      "Create Error",
			title:     "Title",
			content:   "Many samsing words.",
			respError: "internal error",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			postCreaterMock := mocks.NewPostCreater(t)

			if tc.respError == "" || tc.mockError != nil {
				postCreaterMock.On("CreatePost", tc.title, tc.content, mock.Anything).
					Return(1, tc.mockError).Once()
			}

			handler := posts.Create(slogdiscard.NewDiscardLogger(), postCreaterMock)

			input := fmt.Sprintf(`{"title": "%s", "content": "%s"}`, tc.title, tc.content)

			req, err := http.NewRequest(http.MethodPost, "/posts", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp posts.CreatePostResponse

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
