package routes

import (
	"kisahloka_be/controllers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/api/v1/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Selamat Datang di KisahLoka API")
	})

	// Type
	e.GET("/api/v1/type", controllers.GetAllTypes)
	e.GET("/api/v1/type/:type_id", controllers.GetTypeDetail)
	e.POST("/api/v1/type", controllers.CreateType)
	e.PUT("/api/v1/type", controllers.UpdateType)
	e.DELETE("/api/v1/type/:type_id", controllers.DeleteType)

	// Origin
	e.GET("/api/v1/origin", controllers.GetAllOrigins)
	e.GET("/api/v1/origin/:origin_id", controllers.GetOriginDetail)
	e.POST("/api/v1/origin", controllers.CreateOrigin)
	e.PUT("/api/v1/origin", controllers.UpdateOrigin)
	e.DELETE("/api/v1/origin/:origin_id", controllers.DeleteOrigin)

	// Genre
	e.GET("/api/v1/genre", controllers.GetAllGenres)
	e.GET("/api/v1/genre/:genre_id", controllers.GetGenreDetail)
	e.POST("/api/v1/genre", controllers.CreateGenre)
	e.PUT("/api/v1/genre", controllers.UpdateGenre)
	e.DELETE("/api/v1/genre/:genre_id", controllers.DeleteGenre)

	// Role
	e.GET("/api/v1/role", controllers.GetAllRoles)
	e.GET("/api/v1/role/:role_id", controllers.GetRoleDetail)
	e.POST("/api/v1/role", controllers.CreateRole)
	e.PUT("/api/v1/role", controllers.UpdateRole)
	e.DELETE("/api/v1/role/:role_id", controllers.DeleteRole)

	// User
	e.GET("/api/v1/user", controllers.GetAllUsers)
	e.GET("/api/v1/user/:user_id", controllers.GetUserDetail)
	e.GET("/api/v1/user/uid/:uid", controllers.GetUserDetailUID)
	e.POST("/api/v1/user", controllers.CreateUser)
	e.PUT("/api/v1/user", controllers.UpdateUser)
	e.DELETE("/api/v1/user/:user_id", controllers.DeleteUser)

	// Story
	e.GET("/api/v1/story", controllers.GetAllStoriesCompleted)
	e.GET("/api/v1/story_preview", controllers.GetAllStoriesPreview)
	e.GET("/api/v1/story/:story_id", controllers.GetStoryDetail)
	e.GET("/api/v1/story/contents/:story_id", controllers.GetStoryContentOnStory)
	e.POST("/api/v1/story", controllers.CreateStory)
	e.PUT("/api/v1/story", controllers.UpdateStory)
	e.DELETE("/api/v1/story/:story_id", controllers.DeleteStory)
	e.GET("/api/v1/story_recommendation/random/:exclude_story_id", controllers.GetStoriesRecommendationRandom)

	// Bookmark
	e.GET("/api/v1/bookmark", controllers.GetAllBookmarks)
	e.GET("/api/v1/bookmark/user/:user_id", controllers.GetAllBookmarksByUserID)
	e.GET("/api/v1/bookmark/:bookmark_id", controllers.GetBookmarkDetail)
	e.POST("/api/v1/bookmark", controllers.CreateBookmark)
	e.PUT("/api/v1/bookmark", controllers.UpdateBookmark)
	e.DELETE("/api/v1/bookmark/:bookmark_id", controllers.DeleteBookmark)

	// Home
	e.GET("/api/v1/home", controllers.GetHomeData)

	return e
}
