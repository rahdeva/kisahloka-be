// Story Model

package models

import (
	"fmt"
	"kisahloka_be/db"
	"strconv"
	"strings"
	"time"
)

type Story struct {
	StoryID        int                `json:"story_id"`
	TypeID         int                `json:"type_id"`
	TypeName       string             `json:"type_name"`
	OriginID       int                `json:"origin_id"`
	OriginName     string             `json:"origin_name"`
	Title          string             `json:"title"`
	TotalContent   int                `json:"total_content"`
	ReleasedDate   time.Time          `json:"released_date"`
	ThumbnailImage string             `json:"thumbnail_image"`
	ReadCount      int                `json:"read_count"`
	IsHighlighted  int                `json:"is_highligthed"`
	IsFavorited    int                `json:"is_favorited"`
	GenreID        []int              `json:"genre_id"`
	GenreName      []string           `json:"genre_name"`
	Synopsis       string             `json:"synopsis"`
	StoryContent   []StoryContentItem `json:"story_content"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type StoryContentItem struct {
	Order       int    `json:"order"`
	Image       string `json:"image"`
	ContentIndo string `json:"content_indo"`
	ContentEng  string `json:"content_eng"`
}

func GetAllStories(page, pageSize int, keyword string) (Response, error) {
	var res Response
	var arrobj []Story
	var meta Meta

	con := db.CreateCon()

	// Add a WHERE clause to filter stories based on the keyword
	whereClause := ""
	if keyword != "" {
		whereClause = " WHERE title LIKE '%" + keyword + "%'"
	}

	// Count total items in the database
	var totalItems int
	err := con.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM story %s", whereClause)).Scan(&totalItems)
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
			"stories": arrobj,
			"meta":    meta,
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
	sqlStatement := fmt.Sprintf("SELECT s.*, t.type_name, o.origin_name, GROUP_CONCAT(g.genre_id) AS genre_id, GROUP_CONCAT(g.genre_name) AS genre_name FROM story s LEFT JOIN type t ON s.type_id = t.type_id LEFT JOIN origin o ON s.origin_id = o.origin_id LEFT JOIN story_genre sg ON s.story_id = sg.story_id LEFT JOIN genre g ON sg.genre_id = g.genre_id %s GROUP BY s.story_id LIMIT %d OFFSET %d", whereClause, pageSize, offset)
	rows, err := con.Query(sqlStatement)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var obj Story
		var genreIDs, genreNames string
		err := rows.Scan(
			&obj.StoryID,
			&obj.TypeID,
			&obj.OriginID,
			&obj.Title,
			&obj.TotalContent,
			&obj.ReleasedDate,
			&obj.Synopsis,
			&obj.ThumbnailImage,
			&obj.ReadCount,
			&obj.IsHighlighted,
			&obj.IsFavorited,
			&obj.CreatedAt,
			&obj.UpdatedAt,
			&obj.TypeName,
			&obj.OriginName,
			&genreIDs,
			&genreNames,
		)
		if err != nil {
			return res, err
		}

		// Convert time fields to UTC+8 (Asia/Shanghai) before including them in the response
		obj.ReleasedDate = obj.ReleasedDate.In(loc)
		obj.CreatedAt = obj.CreatedAt.In(loc)
		obj.UpdatedAt = obj.UpdatedAt.In(loc)

		// Split genre IDs and names into slices
		obj.GenreID, err = stringsToIntSlice(genreIDs)
		if err != nil {
			return res, err
		}
		obj.GenreName = strings.Split(genreNames, ",")

		// Fetch story content
		content, err := GetStoryContent(obj.StoryID)
		if err != nil {
			return res, err
		}
		obj.StoryContent = content

		arrobj = append(arrobj, obj)

		meta.Limit = pageSize
		meta.Page = page
		meta.TotalPages = calculateTotalPages(totalItems, pageSize)
		meta.TotalItems = totalItems
	}

	res.Data = map[string]interface{}{
		"stories": arrobj,
		"meta":    meta,
	}

	return res, nil
}

func GetStoryContent(storyID int) ([]StoryContentItem, error) {
	var content []StoryContentItem

	con := db.CreateCon()

	sqlStatement := "SELECT `order`, image, content_indo, content_eng FROM story_content WHERE story_id = ? ORDER BY `order`"
	rows, err := con.Query(sqlStatement, storyID)
	if err != nil {
		return content, err
	}
	defer rows.Close()

	for rows.Next() {
		var c StoryContentItem
		err := rows.Scan(
			&c.Order,
			&c.Image,
			&c.ContentIndo,
			&c.ContentEng,
		)
		if err != nil {
			return content, err
		}
		content = append(content, c)
	}

	return content, nil
}

func GetStoryDetail(storyID int) (Response, error) {
	var storyDetail Story
	var res Response

	con := db.CreateCon()

	sqlStatement := `
		SELECT s.story_id, s.type_id, t.type_name, s.origin_id, o.origin_name, 
        s.title, s.total_content, s.released_date, s.thumbnail_image, 
        s.read_count, s.is_highligthed, s.is_favorited, s.synopsis,
		GROUP_CONCAT(sg.genre_id) AS genre_id, GROUP_CONCAT(g.genre_name) AS genre_name
		FROM story s 
		LEFT JOIN type t ON s.type_id = t.type_id 
		LEFT JOIN origin o ON s.origin_id = o.origin_id 
		LEFT JOIN story_genre sg ON s.story_id = sg.story_id 
		LEFT JOIN genre g ON sg.genre_id = g.genre_id 
		WHERE s.story_id = ?
		GROUP BY s.story_id
	`

	row := con.QueryRow(sqlStatement, storyID)

	var genreIDs, genreNames string
	err := row.Scan(
		&storyDetail.StoryID,
		&storyDetail.TypeID,
		&storyDetail.TypeName,
		&storyDetail.OriginID,
		&storyDetail.OriginName,
		&storyDetail.Title,
		&storyDetail.TotalContent,
		&storyDetail.ReleasedDate,
		&storyDetail.ThumbnailImage,
		&storyDetail.ReadCount,
		&storyDetail.IsHighlighted,
		&storyDetail.IsFavorited,
		&storyDetail.Synopsis,
		&genreIDs,
		&genreNames,
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
	storyDetail.ReleasedDate = storyDetail.ReleasedDate.In(loc)

	// Split genre IDs and names into slices
	storyDetail.GenreID = stringsToIntSlice2(genreIDs)
	storyDetail.GenreName = strings.Split(genreNames, ",")

	// Fetch story content
	content, err := GetStoryContent(storyID)
	if err != nil {
		return res, err
	}
	storyDetail.StoryContent = content

	res.Data = map[string]interface{}{
		"story": storyDetail,
	}

	return res, nil
}

func CreateStory(story Story) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "INSERT INTO story (type_id, origin_id, title, total_content, released_date, synopsis, thumbnail_image, read_count, is_highligthed, is_favorited, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	story.CreatedAt = time.Now()
	story.UpdatedAt = time.Now()

	result, err := stmt.Exec(
		story.TypeID,
		story.OriginID,
		story.Title,
		story.TotalContent,
		story.ReleasedDate,
		story.Synopsis,
		story.ThumbnailImage,
		story.ReadCount,
		story.IsHighlighted,
		story.IsFavorited,
		story.CreatedAt,
		story.UpdatedAt,
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
		"created_at": story.CreatedAt.In(loc),
	}

	return res, nil
}

func UpdateStory(storyID int, updateFields map[string]interface{}) (Response, error) {
	var res Response

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Add or update the 'updated_at' field in the updateFields map
	updateFields["updated_at"] = time.Now().In(loc)
	updated_at := updateFields["updated_at"]

	con := db.CreateCon()

	// Construct the SET part of the SQL statement dynamically
	setStatement := "SET "
	values := []interface{}{}
	i := 0

	for fieldName, fieldValue := range updateFields {
		if i > 0 {
			setStatement += ", "
		}
		setStatement += fieldName + " = ?"
		values = append(values, fieldValue)
		i++
	}

	// Construct the final SQL statement
	sqlStatement := "UPDATE story " + setStatement + " WHERE story_id = ?"
	values = append(values, storyID)

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(values...)
	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected": rowsAffected,
		"updated_at":   updated_at,
	}

	return res, nil
}

func DeleteStory(storyID int) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "DELETE FROM story WHERE story_id = ?"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(storyID)

	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected":     rowsAffected,
		"deleted_story_id": storyID,
	}

	return res, err
}

func stringsToIntSlice(s string) ([]int, error) {
	// Pisahkan string menjadi potongan-potongan integer
	parts := strings.Split(s, ",")

	// Buat slice kosong untuk menyimpan hasil konversi
	var result []int

	// Lakukan iterasi pada setiap potongan string
	for _, part := range parts {
		// Konversi potongan string menjadi integer
		num, err := strconv.Atoi(part)
		if err != nil {
			// Jika terjadi kesalahan, kembalikan error
			return nil, err
		}

		// Tambahkan integer ke slice hasil
		result = append(result, num)
	}

	// Kembalikan slice hasil dan tanpa error
	return result, nil
}

func stringsToIntSlice2(s string) []int {
	var result []int
	if s == "" {
		return result
	}
	strArr := strings.Split(s, ",")
	for _, str := range strArr {
		i, err := strconv.Atoi(str)
		if err != nil {
			// Jika terjadi kesalahan dalam mengonversi string menjadi integer, lewati elemen ini
			continue
		}
		result = append(result, i)
	}
	return result
}
