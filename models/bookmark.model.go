package models

import (
	"fmt"
	"kisahloka_be/db"
	"time"
)

type Bookmark struct {
	BookmarkID     int       `json:"bookmark_id"`
	UserID         int       `json:"user_id"`
	UID            string    `json:"uid"`
	StoryID        int       `json:"story_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Title          string    `json:"title"`
	OriginName     string    `json:"origin_name"`
	ThumbnailImage string    `json:"thumbnail_image"`
	TotalContent   int       `json:"total_content"`
}

// GetAllBookmarks retrieves all bookmarks with pagination and optional keyword search
func GetAllBookmarks(page, pageSize int, keyword string) (Response, error) {
	var res Response
	var arrobj []Bookmark
	var meta Meta

	con := db.CreateCon()

	// Add a WHERE clause to filter bookmarks based on the keyword
	whereClause := ""
	if keyword != "" {
		whereClause = " WHERE user_id LIKE '%" + keyword + "%'"
	}

	// Count total items in the database
	var totalItems int
	err := con.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM bookmark %s", whereClause)).Scan(&totalItems)
	if err != nil {
		return res, err
	}

	// If no items are found, return an empty response data object
	if totalItems == 0 {
		meta.Limit = pageSize
		meta.Page = page
		meta.TotalPages = 0
		meta.TotalItems = totalItems

		res.Data = map[string]interface{}{
			"bookmarks": make([]Bookmark, 0), // Empty slice
			"meta":      meta,
		}

		return res, nil
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Calculate the total number of pages
	totalPages := calculateTotalPages(totalItems, pageSize)

	// Check if the requested page is greater than the total number of pages
	if page > totalPages {
		return res, fmt.Errorf("requested page (%d) exceeds total number of pages (%d)", page, totalPages)
	}

	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize
	sqlStatement := fmt.Sprintf("SELECT * FROM bookmark %s LIMIT %d OFFSET %d", whereClause, pageSize, offset)
	rows, err := con.Query(sqlStatement)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var obj Bookmark
		err := rows.Scan(
			&obj.BookmarkID,
			&obj.UserID,
			&obj.StoryID,
			&obj.CreatedAt,
			&obj.UpdatedAt,
		)
		if err != nil {
			return res, err
		}

		// Convert time fields to UTC+8 (Asia/Shanghai) before including them in the response
		obj.CreatedAt = obj.CreatedAt.In(loc)
		obj.UpdatedAt = obj.UpdatedAt.In(loc)

		arrobj = append(arrobj, obj)

		meta.Limit = pageSize
		meta.Page = page
		meta.TotalPages = calculateTotalPages(totalItems, pageSize)
		meta.TotalItems = totalItems
	}

	res.Data = map[string]interface{}{
		"bookmarks": arrobj,
		"meta":      meta,
	}

	return res, nil
}

// GetAllBookmarksByUserID retrieves all bookmarks by user ID with pagination and optional keyword search
func GetAllBookmarksByUserID(userID, page, pageSize int, keyword string) (Response, error) {
	var res Response
	var arrobj []Bookmark
	var meta Meta

	con := db.CreateCon()

	// Add a WHERE clause to filter bookmarks based on the keyword and user ID
	whereClause := fmt.Sprintf("WHERE user_id = %d", userID)
	if keyword != "" {
		whereClause += fmt.Sprintf(" AND story_id LIKE '%%%s%%'", keyword)
	}

	// Count total items in the database for the specific user
	var totalItems int
	err := con.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM bookmark %s", whereClause)).Scan(&totalItems)
	if err != nil {
		return res, err
	}

	// If no items are found, return an empty response data object
	if totalItems == 0 {
		meta.Limit = pageSize
		meta.Page = page
		meta.TotalPages = 0
		meta.TotalItems = totalItems

		res.Data = map[string]interface{}{
			"bookmarks": make([]Bookmark, 0), // Empty slice
			"meta":      meta,
		}

		return res, nil
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Calculate the total number of pages
	totalPages := calculateTotalPages(totalItems, pageSize)

	// Check if the requested page is greater than the total number of pages
	if page > totalPages {
		return res, fmt.Errorf("requested page (%d) exceeds total number of pages (%d)", page, totalPages)
	}

	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize
	sqlStatement := fmt.Sprintf(
		`SELECT 
			bookmark.*, 
			story.title, 
			origin.origin_name, 
			story.thumbnail_image,
			story.total_content
		FROM 
			bookmark 
			INNER JOIN story ON bookmark.story_id = story.story_id 
			INNER JOIN origin ON story.origin_id = origin.origin_id 
			%s 
		LIMIT %d OFFSET %d
		`, whereClause, pageSize, offset)
	rows, err := con.Query(sqlStatement)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var obj Bookmark
		var title, originName, thumbnailImage string
		err := rows.Scan(
			&obj.BookmarkID,
			&obj.UserID,
			&obj.StoryID,
			&obj.CreatedAt,
			&obj.UpdatedAt,
			&title,
			&originName,
			&thumbnailImage,
			&obj.TotalContent,
		)
		if err != nil {
			return res, err
		}

		// Convert time fields to UTC+8 (Asia/Shanghai) before including them in the response
		obj.CreatedAt = obj.CreatedAt.In(loc)
		obj.UpdatedAt = obj.UpdatedAt.In(loc)
		obj.Title = title
		obj.OriginName = originName
		obj.ThumbnailImage = thumbnailImage

		arrobj = append(arrobj, obj)

		meta.Limit = pageSize
		meta.Page = page
		meta.TotalPages = calculateTotalPages(totalItems, pageSize)
		meta.TotalItems = totalItems
	}

	res.Data = map[string]interface{}{
		"bookmarks": arrobj,
		"meta":      meta,
	}

	return res, nil
}

// GetBookmarkDetail retrieves a single bookmark by its ID
func GetBookmarkDetail(bookmarkID int) (Response, error) {
	var bookmarkDetail Bookmark
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT * FROM bookmark WHERE bookmark_id = ?"

	row := con.QueryRow(sqlStatement, bookmarkID)

	err := row.Scan(
		&bookmarkDetail.BookmarkID,
		&bookmarkDetail.UserID,
		&bookmarkDetail.StoryID,
		&bookmarkDetail.CreatedAt,
		&bookmarkDetail.UpdatedAt,
	)

	if err != nil {
		return res, err
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Convert time fields to UTC+8 (Asia/Shanghai) before including them in the response
	bookmarkDetail.CreatedAt = bookmarkDetail.CreatedAt.In(loc)
	bookmarkDetail.UpdatedAt = bookmarkDetail.UpdatedAt.In(loc)

	res.Data = map[string]interface{}{
		"bookmark": bookmarkDetail,
	}

	return res, nil
}

// CreateBookmark creates a new bookmark
func CreateBookmark(userID int, storyID int, uid string) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "INSERT INTO bookmark (user_id, uid, story_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	created_at := time.Now()
	updated_at := time.Now()

	result, err := stmt.Exec(
		userID,
		uid,
		storyID,
		created_at,
		updated_at,
	)

	if err != nil {
		return res, err
	}

	getIDLast, err := result.LastInsertId()

	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"getIDLast":  getIDLast,
		"created_at": created_at.In(loc),
	}

	return res, nil
}

// UpdateBookmark updates an existing bookmark
func UpdateBookmark(bookmarkID, userID, storyID int) (Response, error) {
	var res Response

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	con := db.CreateCon()

	// Construct the SET part of the SQL statement dynamically
	sqlStatement := "UPDATE bookmark SET user_id = ?, story_id = ?, updated_at = ? WHERE bookmark_id = ?"

	// Execute the SQL statement
	result, err := con.Exec(sqlStatement, userID, storyID, time.Now().In(loc), bookmarkID)
	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected": rowsAffected,
	}

	return res, nil
}

// DeleteBookmark deletes a bookmark by its ID
func DeleteBookmark(bookmarkID int) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "DELETE FROM bookmark WHERE bookmark_id = ?"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(bookmarkID)

	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected":        rowsAffected,
		"deleted_bookmark_id": bookmarkID,
	}

	return res, nil
}
