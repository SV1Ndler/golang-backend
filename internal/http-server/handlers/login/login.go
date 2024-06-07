package login

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/slog"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/pkg/models"
)

// URLGetter is an interface for getting url by alias.
//

type loginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	resp.Response
	AccessToken string `json:"access_token"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.40.2 --name=UserGetter
type UserGetter interface {
	GetUserByLoginAndPassword(login string, password string) (models.User, error)
}

func New(log *slog.Logger, userGetter UserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.login.New"
		var jwtSecretKey = []byte("secret") //TODO

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req loginRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty")

			// render.Status(r, http.StatusBadRequest) // ???
			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			// render.Status(r, http.StatusBadRequest) // ???
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		user, err := userGetter.GetUserByLoginAndPassword(req.Login, req.Password)
		if err != nil {
			log.Error("failed to create post", sl.Err(err))

			render.JSON(w, r, resp.Error("login or password is incorrect"))

			return
		}

		log.Info("login succeed")

		payload := jwt.MapClaims{
			"sub": user.Email,
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		}

		// Создаем новый JWT-токен и подписываем его по алгоритму HS256
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

		t, err := token.SignedString(jwtSecretKey)
		if err != nil {
			log.Error("failed to create token", sl.Err(err)) // TODO!!!!
			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		render.JSON(w, r, LoginResponse{
			Response:    resp.OK(),
			AccessToken: t,
		})
	}
}
