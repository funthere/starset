package http

import (
	"net/http"

	"github.com/funthere/starset/productservice/domain"
	"github.com/labstack/echo/v4"
)

type productHandler struct {
	productUsecase domain.ProductUsecase
}

func NewProductHandler(e *echo.Echo, productUc domain.ProductUsecase) {
	handler := productHandler{
		productUsecase: productUc,
	}

	e.POST("/store", handler.Store)

}

func (h *productHandler) Store(c echo.Context) error {
	product := domain.Product{}
	if err := c.Bind(&product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(&product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := h.productUsecase.Store(c.Request().Context(), product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data": &product,
	})
}
