package http

import (
	"net/http"
	"strconv"

	"github.com/funthere/starset/productservice/domain"
	"github.com/funthere/starset/productservice/middlewares"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type productHandler struct {
	productUsecase domain.ProductUsecase
}

func NewProductHandler(e *echo.Echo, productUc domain.ProductUsecase) {
	handler := &productHandler{productUc}

	router := e.Group("product")
	router.Use(middlewares.Authentication)

	router.POST("", handler.Store)
	router.GET("/:id", handler.Get)
	router.GET("", handler.Fetch)

}

func (h *productHandler) Store(c echo.Context) error {
	product := domain.Product{}
	if err := c.Bind(&product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(&product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	userData := c.Get("userData").(jwt.MapClaims)
	userID := uint32(userData["id"].(float64))
	product.OwnerID = userID
	product.Owner.ID = userID

	if err := h.productUsecase.Store(c.Request().Context(), &product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"data": &product,
	})
}

func (h *productHandler) Get(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	res, err := h.productUsecase.GetById(c.Request().Context(), uint32(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h *productHandler) Fetch(c echo.Context) error {
	filter := domain.Filter{
		Search: c.QueryParam("search"),
	}

	res, err := h.productUsecase.Fetch(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
