package middle

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/mlbautomation/ProyectoEMLB/infrastructure/handlers/response"
	"github.com/mlbautomation/ProyectoEMLB/model"
)

type AuthMiddleware struct {
	responser response.API
}

func New() *AuthMiddleware {
	return &AuthMiddleware{}
}

// valida si el token es válido, y devuelve los customclaims
func (am *AuthMiddleware) IsValid(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obtiene el token de body
		token, err := getTokenFromRequest(c.Request())
		if err != nil {
			return am.responser.Error(c, "middle-AuthMiddleware-IsValid-getTokenFromRequest(c.Request())", err)
		}
		// valida el token
		isValid, claims := am.validate(token)
		if !isValid {
			err = errors.New("the token is not valid")
			return am.responser.Error(c, "middle-AuthMiddleware-IsValid-am.validate(token)", err)
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("isAdmin", claims.IsAdmin)

		return next(c)
	}
}

func (am *AuthMiddleware) IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		isAdmin, ok := c.Get("isAdmin").(bool)
		if !isAdmin || !ok {
			err := errors.New("you are not admin")
			return am.responser.Error(c, "middle-AuthMiddleware-IsAdmin-c.Get('isAdmin').(bool)", err)
		}

		return next(c)
	}
}

// valida si el token es válido, y devuelve los customclaims
func (am *AuthMiddleware) validate(token string) (bool, model.JWTCustomClaims) {
	//super duda: ¿el SECRET KEY entonces sirve tambien para desencriptar!?
	claims, err := jwt.ParseWithClaims(token, &model.JWTCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		log.Println(token)
		log.Println(os.Getenv("JWT_SECRET_KEY"))
		log.Println(err)
		return false, model.JWTCustomClaims{}
	}

	data, ok := claims.Claims.(*model.JWTCustomClaims)
	if !ok {
		log.Println("is not a jwtcustomclaims")
		return false, model.JWTCustomClaims{}
	}

	return true, *data
}

// obtiene el token de body
func getTokenFromRequest(r *http.Request) (string, error) {
	data := r.Header.Get("Authorization")
	if data == "" {
		return "", errors.New("authorization in headers is empty")
	}

	if strings.HasPrefix(data, "Bearer") {
		return data[7:], nil
	}

	return data, nil
}
