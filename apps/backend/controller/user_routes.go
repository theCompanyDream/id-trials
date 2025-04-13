package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	model "github.com/theCompanyDream/id-trials/apps/backend/models"
	db "github.com/theCompanyDream/id-trials/apps/backend/repository"
)

// GetUser godoc
// @Summary Get a single user
// @Description Get a user by their ID or username
// @Tags user
// @Accept json
// @Produce json
// @Param id path string false "User ID"
// @Param user_name path string false "Username"
// @Success 302 {object} models.UserUlid "User Found"
// @Failure 400 {object} object "Bad Request"
// @Router /user/{id} [get]
func GetUser(c echo.Context) error {
	// Extract the user ID from the URL and query the database
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusNotFound, errors.New("id not applicable there"))
	}
	user, err := db.GetUser(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

// GetUsers godoc
// @Summary Get multiple users
// @Description Get a list of users, with optional search, pagination, and limit
// @Tags user
// @Accept json
// @Produce json
// @Param search query string false "Search Term"
// @Param limit query int false "Limit"
// @Param page query int false "Page Number"
// @Success 302 {object} []models.UserUlidPaging "Users Found"
// @Failure 400 {object} object "Bad Request"
// @Router /users [get]
func GetUsers(c echo.Context) error {
	// Extract the user ID from the URL and query the database
	var page, limit int
	search := c.QueryParam("search")
	if queryLimit := c.QueryParam("limit"); queryLimit != "" {
		i, _ := strconv.Atoi(queryLimit)
		limit = i
	} else {
		limit = 25
	}
	if queryPage := c.QueryParam("page"); queryPage != "" {
		i, _ := strconv.Atoi(queryPage)
		page = i
	} else {
		page = 1
	}
	users, error := db.GetUsers(search, page, limit, c)
	if error != nil {
		return error
	}
	c.Logger().Info(users)
	return c.JSON(http.StatusOK, users)
}

// CreateUser godoc
// @Summary Create a user
// @Description Create a new user with the provided information
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.UserInput true "User object"
// @Success 201 {object} models.UserUlid "User Created"
// @Failure 400 {object} object "Bad Request"
// @Router /user [post]
func CreateUser(c echo.Context) error {
	// Parse user details from the request body and insert into the database
	request := model.UserInput{}
	err := c.Bind(&request)
	if err != nil {
		return err
	}
	err = validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusUnprocessableEntity, validationErrorsToMap(validationErrors))
	}
	dto := model.InputToDTO(request)
	user, error := db.CreateUser(*dto)
	if error != nil {
		return error
	}
	return c.JSON(http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user's information by their ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UserInput true "User object"
// @Success 200 {object} models.UserUlid "User Updated"
// @Failure 400 {object} object "Bad Request"
// @Router /user/{id} [put]
func UpdateUser(c echo.Context) error {
	// Parse user details from the request body and insert into the database
	// request := checkConstraints(c)
	request := model.UserInput{}
	err := c.Bind(&request)
	if err != nil {
		return err
	}
	err = validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusUnprocessableEntity, validationErrorsToMap(validationErrors))
	}
	if id := c.Param("id"); id != "" {
		request.Id = &id
	}
	dto := model.InputToDTO(request)
	user, error := db.UpdateUser(*dto)
	if error != nil {
		return error
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {string} string "User Deleted"
// @Failure 400 {object} object "Bad Request"
// @Router /user/{id} [delete]
func DeleteUser(c echo.Context) error {
	// Parse user details from the request body and insert into the database
	id := c.Param("id")
	if id == "" {
		return errors.New("id must not be null")
	}
	err := db.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
