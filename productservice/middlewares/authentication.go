package middlewares

import (
	"net/http"

	"github.com/funthere/starset/productservice/helper"
	"github.com/labstack/echo/v4"
)

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		verifyToken, err := helper.VerifyToken(c)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, err)

		}

		c.Set("userData", verifyToken)
		return next(c)
	}
}
