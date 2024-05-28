// Story Controller

package controllers

import (
	"kisahloka_be/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllStoriesCompleted(c echo.Context) error {
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

	result, err := models.GetAllStoriesCompleted(page, pageSize, keyword)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func GetAllStoriesPreview(c echo.Context) error {
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

	result, err := models.GetAllStoriesPreview(page, pageSize, keyword)
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

	userIDParam := c.QueryParam("user_id")
	var userID *int
	if userIDParam != "" {
		parsedUserID, err := strconv.Atoi(userIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id"})
		}
		userID = &parsedUserID
	}

	uidParam := c.QueryParam("uid")
	var uid *string
	if uidParam != "" {
		uid = &uidParam
	}

	storyDetail, err := models.GetStoryDetail(storyID, userID, uid)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, storyDetail)
}

func GetStoryContentOnStory(c echo.Context) error {
	storyID, err := strconv.Atoi(c.Param("story_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid story_id"})
	}

	storyDetail, err := models.GetStoryContentOnStory(storyID)
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

func GetStoriesRecommendationRandom(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 4 // Default limit
	}

	excludeStoryID, err := strconv.Atoi(c.Param("exclude_story_id"))
	if err != nil {
		excludeStoryID = 0 // Default exclude story ID
	}

	result, err := models.GetStoriesRecommendationRandom(limit, excludeStoryID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}
