package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/MoonSHRD/sonis/controllers"
	"github.com/MoonSHRD/sonis/utils"
	"github.com/dgrijalva/jwt-go"

	"github.com/MoonSHRD/logger"
	"github.com/MoonSHRD/sonis/app"
	"github.com/labstack/echo/v4"
)

const (
	ValidateTokenEndpoint = "/api/v1/validateAccessToken"
)

type AuthMiddleware struct {
	app *app.App
}

func NewAuthMiddleware(a *app.App) *AuthMiddleware {
	return &AuthMiddleware{
		app: a,
	}
}

func (am *AuthMiddleware) ProcessRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authServerValidationEndpoint := path.Join(am.app.Config.AuthServerURL, ValidateTokenEndpoint)

		// parse access token from http header (Bearer Auth)
		accessToken := c.Request().Header.Get("Authorization")
		if len(accessToken) == 0 {
			err := fmt.Errorf("unauthorized")
			logger.Warningf("Unauthorized access on route %s", c.Request().URL.Path)
			controllers.ReturnHTTPError(c, err, http.StatusForbidden)
			return err
		}
		splitToken := strings.Split(accessToken, " ")
		accessToken = splitToken[1]

		validateAccessTokenBodyRequest, err := json.Marshal(map[string]string{
			"accessToken": accessToken,
		})
		resp, err := http.Post(authServerValidationEndpoint, "application/json", bytes.NewBuffer(validateAccessTokenBodyRequest))
		if err != nil {
			logger.Error(err)
			controllers.ReturnHTTPError(c, err, http.StatusInternalServerError)
			return err
		}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Error(err)
				controllers.ReturnHTTPError(c, err, http.StatusInternalServerError)
				return err
			}
			resp.Body.Close()
			var authServerResponse struct {
				Valid bool `json:"valid"`
			}
			err = json.Unmarshal([]byte(bodyBytes), &authServerResponse)
			if err != nil {
				logger.Error(err)
				controllers.ReturnHTTPError(c, err, http.StatusInternalServerError)
				return err
			}
			if authServerResponse.Valid {
				token, err := jwt.ParseWithClaims(accessToken, &utils.JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(am.app.Config.JWT.SigningKey), nil
				})

				if err != nil {
					logger.Error(err)
					controllers.ReturnHTTPError(c, err, http.StatusInternalServerError)
					return err
				}

				c.Set("claims", token)
				next(c)
			} else {
				err := fmt.Errorf("incorrect")
				logger.Warningf("Unauthorized access on route %s", c.Request().URL.Path)
				controllers.ReturnHTTPError(c, err, http.StatusForbidden)
				return err
			}
		}
		err = fmt.Errorf("Incorrect HTTP code from auth server - %d", resp.StatusCode)
		logger.Warningf(err.Error())
		controllers.ReturnHTTPError(c, err, http.StatusInternalServerError)
		return err
	}
}
