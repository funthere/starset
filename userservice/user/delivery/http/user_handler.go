package http

import (
	"net/http"

	"github.com/funthere/starset/userservice/domain"
	"github.com/funthere/starset/userservice/helpers"
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
