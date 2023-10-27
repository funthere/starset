package helper

import (
	"errors"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var secretKey = os.Getenv("JWT_SECRET")

func GenerateToken(id uint32, email string) string {
	// menyimpan data user
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
	}

	// pilih metode enkripsi
	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// parsing token menjadi string
	signedToken, _ := parseToken.SignedString([]byte(secretKey))

	return signedToken
}

func VerifyToken(c echo.Context) (interface{}, error) {
	errResponse := errors.New("sign in to proceed")
	headerToken := c.Request().Header.Get("Authorization")
	bearer := strings.HasPrefix(headerToken, "Bearer")

	if !bearer {
		return nil, errResponse
	}

	stringToken := strings.Split(headerToken, " ")[1]

	token, _ := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResponse
		}
		return []byte(secretKey), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, errResponse
	}

	return token.Claims.(jwt.MapClaims), nil
}
