package http

import (
	"net/http"
	"strconv"
	"strings"

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
	router.GET("/fetch", handler.FetchByIds)

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

func (h *productHandler) FetchByIds(c echo.Context) error {
	var (
		strIds   = strings.Split(c.QueryParam("ids"), ",")
		response = map[uint32]domain.Product{}
	)

	if len(strIds) >= 1 && strIds[0] != "" {
		ids := []uint32{}
		for _, val := range strIds {
			ui64, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			ids = append(ids, uint32(ui64))
		}
		filter := domain.Filter{
			IDs: ids,
		}

		products, err := h.productUsecase.Fetch(c.Request().Context(), filter)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		for _, val := range products {
			response[val.ID] = val
		}
	}

	return c.JSON(http.StatusOK, response)
}
