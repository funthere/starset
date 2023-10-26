package http

import (
	"net/http"
	"os"
	"time"

	"github.com/funthere/starset/userservice/domain"
	"github.com/funthere/starset/userservice/helpers"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(e *echo.Echo, userUc domain.UserUsecase) {
	handler := UserHandler{
		userUsecase: userUc,
	}

	e.POST("/register", handler.RegisterUser)
	e.POST("/login", handler.LoginUser)
}

func (h UserHandler) RegisterUser(c echo.Context) error {
	user := domain.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "invalid json payload"})
	}

	if err := c.Validate(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := h.userUsecase.Register(c.Request().Context(), user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"data": &user,
	})
}

func (h UserHandler) LoginUser(c echo.Context) error {
	user := domain.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "invalid json payload"})
	}

	password := user.Password
	if err := h.userUsecase.Login(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// compare password
	comparePass := helpers.ComparePassword([]byte(user.Password), []byte(password))
	if !comparePass {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.ErrResponse{Message: "invalid password"})
	}

	// generate jwt token
	jwtToken, err := GenerateJwtToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"token": jwtToken,
	})
}

var secretKey = os.Getenv("JWT_SECRET")

func GenerateJwtToken(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"address": user.Address,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
