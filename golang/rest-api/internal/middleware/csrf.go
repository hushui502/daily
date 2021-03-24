package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"rest-api/pkg/csrf"
	"rest-api/pkg/httpErrors"
	"rest-api/pkg/utils"
)

func (mw *MiddlewareManager) CSRF(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !mw.cfg.Server.CSRF {
			return next(c)
		}

		token := c.Request().Header.Get(csrf.CSRFHeader)
		if token == "" {
			mw.logger.Errorf("CSRF Middleware get CSRF header, Token: %s, Error: %s, RequestId: %s",
				token,
				"empty CSRF token",
				utils.GetRequestID(c),
			)
			return c.JSON(http.StatusForbidden, httpErrors.NewRestError(http.StatusForbidden, "Invalid CSRF Token", "no CSRF Token"))
		}

		sid, ok := c.Get("sid").(string)
		if !csrf.ValidateToken(token, sid, mw.logger) || !ok {
			mw.logger.Errorf("CSRF Middleware csrf.ValidateToken Token: %s, Error: %s, RequestId: %s",
				token,
				"empty token",
				utils.GetRequestID(c),
			)
			return c.JSON(http.StatusForbidden, httpErrors.NewRestError(http.StatusForbidden, "Invalid CSRF Token", "no CSRF Token"))
		}

		return next(c)
	}
}
