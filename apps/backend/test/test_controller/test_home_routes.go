package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/theCompanyDream/id-trials/apps/backend/controller"
)

func TestHome(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := controller.Home(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Server is up and running")

	// Or if you want to check the exact JSON:
	expected := `{"data":"Server is up and running"}`
	assert.JSONEq(t, expected, rec.Body.String())
}
