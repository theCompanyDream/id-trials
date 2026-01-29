package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
)

func HttpErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal server error"

	// Check if it's an Echo HTTPError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if msg, ok := he.Message.(string); ok {
			message = msg
		}
	}

	// ✅ Use your zerolog Logger instead of gommon/log
	Logger.LogError().
		Str("method", c.Request().Method).
		Str("uri", c.Request().URL.Path).
		Str("error", err.Error()).
		Int("status_code", code).
		Str("stack_trace", string(debug.Stack())).
		Msg("HTTP Error")

	// Respond to client
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			c.NoContent(code)
		} else {
			c.JSON(code, map[string]string{
				"error": message,
			})
		}
	}
}
