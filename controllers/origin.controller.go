// Origin controller
package controllers

import (
	"kisahloka_be/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllOrigins(c echo.Context) error {
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

	result, err := models.GetAllOrigins(page, pageSize, keyword)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func GetOriginDetail(c echo.Context) error {
	originID, err := strconv.Atoi(c.Param("origin_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid origin_id"})
	}

	originDetail, err := models.GetOriginDetail(originID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, originDetail)
}

func CreateOrigin(c echo.Context) error {
	var originObj models.Origin

	// Parse the request body to populate the origin struct
	if err := c.Bind(&originObj); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid request body",
			},
		)
	}

	// Call the CreateOrigin function from the models package
	result, err := models.CreateOrigin(originObj.OriginName)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateOrigin(c echo.Context) error {
	// Parse the request body to get the update data
	var updateFields map[string]interface{}
	if err := c.Bind(&updateFields); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid request body"},
		)
	}

	// Extract the ID from the update data
	originID, ok := updateFields["origin_id"].(float64)
	if !ok {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid origin_id format"},
		)
	}

	// Convert id to integer
	convID := int(originID)

	// Remove id from the updateFields map before passing it to the model
	delete(updateFields, "origin_id")

	// Call the UpdateOrigin function from the models package
	result, err := models.UpdateOrigin(convID, updateFields)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteOrigin(c echo.Context) error {
	originID, err := strconv.Atoi(c.Param("origin_id"))
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	result, err := models.DeleteOrigin(originID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}
