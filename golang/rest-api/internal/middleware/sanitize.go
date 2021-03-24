package middleware

import (
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"rest-api/pkg/sanitize"
)

func (mw *MiddlewareManager) Sanitize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		defer c.Request().Body.Close()

		sanBody, err := sanitize.SanitizeJSON(body)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		c.Set("body", sanBody)
		return next(c)
	}
}
