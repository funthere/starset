package http

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/funthere/starset/userservice/domain"
	"github.com/funthere/starset/userservice/helpers"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(e *echo.Echo, userUc domain.UserUsecase) {
	handler := userHandler{
		userUsecase: userUc,
	}

	e.POST("/register", handler.RegisterUser)
	e.POST("/login", handler.LoginUser)
	e.GET("/users", handler.FetchUser)
}

func (h userHandler) RegisterUser(c echo.Context) error {
	user := new(domain.User)
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "invalid json payload"})
	}

	if err := c.Validate(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := h.userUsecase.Register(c.Request().Context(), user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"data": user,
	})
}

func (h userHandler) LoginUser(c echo.Context) error {
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

func (h userHandler) FetchUser(c echo.Context) error {
	var (
		strIds   = strings.Split(c.QueryParam("ids"), ",")
		response = map[uint32]domain.User{}
	)

	if len(strIds) >= 1 && strIds[0] != "" {
		var ids []uint32
		for _, val := range strIds {
			ui64, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
			}
			ids = append(ids, uint32(ui64))
		}

		users, err := h.userUsecase.FetchUsersByIds(c.Request().Context(), ids)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		for _, val := range users {
			response[val.ID] = val
		}
	}

	return c.JSON(http.StatusOK, response)
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
