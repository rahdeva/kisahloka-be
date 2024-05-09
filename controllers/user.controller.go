// User Controller
package controllers

import (
	"kisahloka_be/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetAllUsers returns all users with pagination and optional keyword search
func GetAllUsers(c echo.Context) error {
	// Get query parameters for pagination
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	keyword := c.QueryParam("keyword")

	result, err := models.GetAllUsers(page, pageSize, keyword)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

// GetUserDetail returns details of a specific user by its ID
func GetUserDetail(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id"})
	}

	userDetail, err := models.GetUserDetail(userID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, userDetail)
}

// CreateUser creates a new user with the provided data
func CreateUser(c echo.Context) error {
	var userObj models.User

	// Parse the request body to populate the user struct
	if err := c.Bind(&userObj); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid request body",
			},
		)
	}

	// Call the CreateUser function from the models package
	result, err := models.CreateUser(userObj.Email, userObj.Name, userObj.RoleID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

// UpdateUser updates an existing user with the provided ID and fields
func UpdateUser(c echo.Context) error {
	// Parse the request body to get the update data
	var updateFields map[string]interface{}
	if err := c.Bind(&updateFields); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid request body"},
		)
	}

	// Extract the ID from the update data
	userID, ok := updateFields["user_id"].(float64)
	if !ok {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid user_id format"},
		)
	}

	// Convert id to integer
	convID := int(userID)

	// Remove id from the updateFields map before passing it to the model
	delete(updateFields, "user_id")

	// Call the UpdateUser function from the models package
	result, err := models.UpdateUser(convID, updateFields)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

// DeleteUser deletes a user with the provided ID
func DeleteUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	result, err := models.DeleteUser(userID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}
