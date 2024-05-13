package models

import (
	"kisahloka_be/db"
	"time"
)

type Home struct {
	HighlightStories []StoryHome     `json:"highlight_stories"`
	FavoriteStories  []StoryHome     `json:"favorite_stories"`
	StoryTypes       []StoryTypeHome `json:"story_types"`
}

type StoryHome struct {
	StoryID        int       `json:"story_id"`
	TypeID         int       `json:"type_id"`
	TypeName       string    `json:"type_name"`
	OriginID       int       `json:"origin_id"`
	OriginName     string    `json:"origin_name"`
	Title          string    `json:"title"`
	ThumbnailImage string    `json:"thumbnail_image"`
	IsHighlighted  int       `json:"is_highligthed"`
	IsFavorited    int       `json:"is_favorited"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type StoryTypeHome struct {
	TypeID   int    `json:"type_id"`
	TypeName string `json:"type_name"`
}

func GetHomeData() (Response, error) {
	var res Response
	var homeData Home

	// Fetching highlighted stories
	highlightedStories, err := getStoriesByHighlight(1)
	if err != nil {
		res.Error = err.Error()
		return res, err
	}
	homeData.HighlightStories = highlightedStories

	// Fetching favorite stories
	favoriteStories, err := getStoriesByFavorite(1)
	if err != nil {
		res.Error = err.Error()
		return res, err
	}
	homeData.FavoriteStories = favoriteStories

	// Fetching all story types
	storyTypes, err := getAllStoryTypes()
	if err != nil {
		res.Error = err.Error()
		return res, err
	}
	homeData.StoryTypes = storyTypes

	// Set the data in the Response struct
	res.Data = struct {
		HighlightStories []StoryHome     `json:"highlight_stories"`
		FavoriteStories  []StoryHome     `json:"favorite_stories"`
		StoryTypes       []StoryTypeHome `json:"story_types"`
	}{
		HighlightStories: homeData.HighlightStories,
		FavoriteStories:  homeData.FavoriteStories,
		StoryTypes:       homeData.StoryTypes,
	}

	return res, nil
}

func getStoriesByHighlight(highlight int) ([]StoryHome, error) {
	var stories []StoryHome

	db := db.CreateCon()

	rows, err := db.Query("SELECT s.story_id, s.type_id, t.type_name, s.origin_id, o.origin_name, s.title, s.thumbnail_image, s.is_highligthed, s.is_favorited, s.created_at, s.updated_at FROM story s JOIN type t ON s.type_id = t.type_id JOIN origin o ON s.origin_id = o.origin_id WHERE s.is_highligthed = ?", highlight)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var story StoryHome
		err := rows.Scan(&story.StoryID, &story.TypeID, &story.TypeName, &story.OriginID, &story.OriginName, &story.Title, &story.ThumbnailImage, &story.IsHighlighted, &story.IsFavorited, &story.CreatedAt, &story.UpdatedAt)
		if err != nil {
			return nil, err
		}
		stories = append(stories, story)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stories, nil
}

func getStoriesByFavorite(favorite int) ([]StoryHome, error) {
	var stories []StoryHome

	db := db.CreateCon()

	rows, err := db.Query("SELECT s.story_id, s.type_id, t.type_name, s.origin_id, o.origin_name, s.title, s.thumbnail_image, s.is_highligthed, s.is_favorited, s.created_at, s.updated_at FROM story s JOIN type t ON s.type_id = t.type_id JOIN origin o ON s.origin_id = o.origin_id WHERE s.is_favorited = ?", favorite)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var story StoryHome
		err := rows.Scan(&story.StoryID, &story.TypeID, &story.TypeName, &story.OriginID, &story.OriginName, &story.Title, &story.ThumbnailImage, &story.IsHighlighted, &story.IsFavorited, &story.CreatedAt, &story.UpdatedAt)
		if err != nil {
			return nil, err
		}
		stories = append(stories, story)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stories, nil
}

func getAllStoryTypes() ([]StoryTypeHome, error) {
	var storyTypes []StoryTypeHome

	db := db.CreateCon()

	rows, err := db.Query("SELECT type_id, type_name FROM type")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var storyType StoryTypeHome
		err := rows.Scan(&storyType.TypeID, &storyType.TypeName)
		if err != nil {
			return nil, err
		}
		storyTypes = append(storyTypes, storyType)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return storyTypes, nil
}
