package http

import (
	"net/http"
	"strconv"

	"github.com/funthere/starset/orderservice/domain"
	"github.com/funthere/starset/orderservice/middlewares"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type orderHandler struct {
	orderUsecase domain.OrderUsecase
}

func NewOrderHandler(e *echo.Echo, orderUc domain.OrderUsecase) {
	handler := &orderHandler{orderUc}

	router := e.Group("order")
	router.Use(middlewares.Authentication)

	router.POST("", handler.Store)
	router.GET("", handler.Fetch)
	router.PATCH("/:id", handler.PatchStatus)
}

func (h *orderHandler) Store(c echo.Context) error {
	order := domain.Order{}
	if err := c.Bind(&order); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(&order); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	userData := c.Get("userData").(jwt.MapClaims)
	userID := userData["id"].(int64)
	order.BuyerID = userID

	if err := h.orderUsecase.Store(c.Request().Context(), &order); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, &order)
}

func (h *orderHandler) Fetch(c echo.Context) error {
	userData := c.Get("userData").(jwt.MapClaims)
	userID := uint32(userData["id"].(float64))
	filter := domain.Filter{
		BuyerID: int64(userID),
	}

	res, err := h.orderUsecase.Fetch(c.Request().Context(), filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *orderHandler) PatchStatus(c echo.Context) error {
	type patchRequest struct {
		Status string `json:"status" validate:"required"`
	}

	var req patchRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if req.Status != domain.StatusAccepted || req.Status != domain.StatusRejected {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid status")
	}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.orderUsecase.PatchStatus(c.Request().Context(), id, req.Status); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, nil)
}
