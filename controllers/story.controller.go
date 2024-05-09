// Story Controller

package controllers

import (
	"kisahloka_be/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllStories(c echo.Context) error {
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

	result, err := models.GetAllStories(page, pageSize, keyword)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func GetStoryDetail(c echo.Context) error {
	storyID, err := strconv.Atoi(c.Param("story_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid story_id"})
	}

	storyDetail, err := models.GetStoryDetail(storyID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, storyDetail)
}

func CreateStory(c echo.Context) error {
	var storyObj models.Story

	// Parse the request body to populate the story struct
	if err := c.Bind(&storyObj); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid request body",
			},
		)
	}

	// Call the CreateStory function from the models package
	result, err := models.CreateStory(storyObj)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateStory(c echo.Context) error {
	// Parse the request body to get the update data
	var updateFields map[string]interface{}
	if err := c.Bind(&updateFields); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid request body"},
		)
	}

	// Extract the ID from the update data
	storyID, ok := updateFields["story_id"].(float64)
	if !ok {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid story_id format"},
		)
	}

	// Convert id to integer
	convID := int(storyID)

	// Remove id from the updateFields map before passing it to the model
	delete(updateFields, "story_id")

	// Call the UpdateStory function from the models package
	result, err := models.UpdateStory(convID, updateFields)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteStory(c echo.Context) error {
	storyID, err := strconv.Atoi(c.Param("story_id"))
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	result, err := models.DeleteStory(storyID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}
