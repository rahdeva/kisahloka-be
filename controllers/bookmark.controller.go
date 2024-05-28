package controllers

import (
	"kisahloka_be/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetAllBookmarks retrieves all bookmarks with pagination and optional keyword search
func GetAllBookmarks(c echo.Context) error {
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

	result, err := models.GetAllBookmarks(page, pageSize, keyword)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

// GetAllBookmarksByUserID retrieves all bookmarks by user ID with pagination and optional keyword search
func GetAllBookmarksByUserID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id"})
	}

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

	result, err := models.GetAllBookmarksByUserID(userID, page, pageSize, keyword)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

// GetBookmarkDetail retrieves a single bookmark by its ID
func GetBookmarkDetail(c echo.Context) error {
	bookmarkID, err := strconv.Atoi(c.Param("bookmark_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid bookmark_id"})
	}

	bookmarkDetail, err := models.GetBookmarkDetail(bookmarkID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, bookmarkDetail)
}

// CreateBookmark creates a new bookmark
func CreateBookmark(c echo.Context) error {
	var bookmarkData struct {
		UserID  int    `json:"user_id"`
		StoryID int    `json:"story_id"`
		UID     string `json:"uid"`
	}

	// Parse the request body to populate bookmarkData struct
	if err := c.Bind(&bookmarkData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	result, err := models.CreateBookmark(bookmarkData.UserID, bookmarkData.StoryID, bookmarkData.UID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// UpdateBookmark updates an existing bookmark
func UpdateBookmark(c echo.Context) error {
	bookmarkID, err := strconv.Atoi(c.Param("bookmark_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid bookmark_id"})
	}

	userID, err := strconv.Atoi(c.FormValue("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id"})
	}

	storyID, err := strconv.Atoi(c.FormValue("story_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid story_id"})
	}

	result, err := models.UpdateBookmark(bookmarkID, userID, storyID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

// DeleteBookmark deletes a bookmark by its ID
func DeleteBookmark(c echo.Context) error {
	bookmarkID, err := strconv.Atoi(c.Param("bookmark_id"))
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	result, err := models.DeleteBookmark(bookmarkID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}
