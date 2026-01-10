package controller

import (
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func HttpErrorHandler(err error, c echo.Context) {
	// Log the error in a more structured format
	log.Errorf("Error: %s, Stack Trace: %s", err.Error(), debug.Stack())

	// Respond with a generic error message to the client
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			c.NoContent(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		}
	}
}
