package server

import (
	"net/http"

	"github.com/ikedam/pictmanager/pkg/rfc7807"
	"github.com/labstack/echo/v4"
)

func HandleRFC7807ErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(rawError error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		err := rfc7807.As(rawError)
		if err == nil {
			e.DefaultHTTPErrorHandler(rawError, c)
			return
		}

		if c.Request().Method == http.MethodHead {
			err := c.NoContent(err.Status)
			if err != nil {
				e.Logger.Error(err)
			}
		} else {
			err := c.JSON(err.Status, err)
			if err != nil {
				e.Logger.Error(err)
			}
		}
	}
}
