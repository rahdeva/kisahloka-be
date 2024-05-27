// User model
package models

import (
	"fmt"
	"kisahloka_be/db"
	"reflect"
	"time"
)

type User struct {
	UserID    int       `json:"user_id"`
	UID       string    `json:"uid"`
	RoleID    int       `json:"role_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	BirthDate time.Time `json:"birth_date"`
	Gender    string    `json:"gender"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAllUsers(page, pageSize int, keyword string) (Response, error) {
	var res Response
	var arrobj reflect.Value
	var meta Meta

	con := db.CreateCon()

	whereClause := ""
	if keyword != "" {
		whereClause = " WHERE email LIKE '%" + keyword + "%' OR name LIKE '%" + keyword + "%'"
	}

	var totalItems int
	err := con.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM user %s", whereClause)).Scan(&totalItems)
	if err != nil {
		return res, err
	}

	if totalItems == 0 {
		meta.Limit = pageSize
		meta.Page = page
		meta.TotalPages = 0
		meta.TotalItems = totalItems

		res.Data = map[string]interface{}{
			"users": make([]interface{}, 0),
			"meta":  meta,
		}

		return res, nil
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	totalPages := calculateTotalPages(totalItems, pageSize)

	if page > totalPages {
		return res, fmt.Errorf("requested page (%d) exceeds total number of pages (%d)", page, totalPages)
	}

	offset := (page - 1) * pageSize
	sqlStatement := fmt.Sprintf("SELECT * FROM user %s LIMIT %d OFFSET %d", whereClause, pageSize, offset)
	rows, err := con.Query(sqlStatement)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var obj User
		err := rows.Scan(
			&obj.UserID,
			&obj.UID,
			&obj.RoleID,
			&obj.Email,
			&obj.Name,
			&obj.BirthDate,
			&obj.Gender,
			&obj.CreatedAt,
			&obj.UpdatedAt,
		)
		if err != nil {
			return res, err
		}

		obj.CreatedAt = obj.CreatedAt.In(loc)
		obj.UpdatedAt = obj.UpdatedAt.In(loc)

		if !arrobj.IsValid() {
			arrobj = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(obj)), 0, 0)
		}

		arrobj = reflect.Append(arrobj, reflect.ValueOf(obj))

		meta.Limit = pageSize
		meta.Page = page
		meta.TotalPages = calculateTotalPages(totalItems, pageSize)
		meta.TotalItems = totalItems
	}

	res.Data = map[string]interface{}{
		"users": arrobj.Interface(),
		"meta":  meta,
	}

	return res, nil
}

func GetUserDetail(userID int) (Response, error) {
	var userDetail User
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT * FROM user WHERE user_id = ?"

	row := con.QueryRow(sqlStatement, userID)

	err := row.Scan(
		&userDetail.UserID,
		&userDetail.UID,
		&userDetail.RoleID,
		&userDetail.Email,
		&userDetail.Name,
		&userDetail.BirthDate,
		&userDetail.Gender,
		&userDetail.CreatedAt,
		&userDetail.UpdatedAt,
	)

	if err != nil {
		return res, err
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	userDetail.CreatedAt = userDetail.CreatedAt.In(loc)
	userDetail.UpdatedAt = userDetail.UpdatedAt.In(loc)

	res.Data = map[string]interface{}{
		"user": userDetail,
	}

	return res, nil
}

func GetUserDetailUID(uid string) (Response, error) {
	var userDetail User
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT * FROM user WHERE uid = ?"

	row := con.QueryRow(sqlStatement, uid)

	err := row.Scan(
		&userDetail.UserID,
		&userDetail.UID,
		&userDetail.RoleID,
		&userDetail.Email,
		&userDetail.Name,
		&userDetail.BirthDate,
		&userDetail.Gender,
		&userDetail.CreatedAt,
		&userDetail.UpdatedAt,
	)

	if err != nil {
		return res, err
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	userDetail.CreatedAt = userDetail.CreatedAt.In(loc)
	userDetail.UpdatedAt = userDetail.UpdatedAt.In(loc)

	res.Data = map[string]interface{}{
		"user": userDetail,
	}

	return res, nil
}

func CreateUser(uid string, roleID int, email, name, gender string, birthDate time.Time) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "INSERT INTO user (uid, role_id, email, name, birth_date, gender, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	created_at := time.Now()
	updated_at := time.Now()

	result, err := stmt.Exec(
		uid,
		roleID,
		email,
		name,
		birthDate,
		gender,
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
		"userID":     getIDLast,
		"created_at": created_at.In(loc),
	}

	return res, nil
}

// UpdateUser updates an existing user with the provided ID and fields
func UpdateUser(userID int, updateFields map[string]interface{}) (Response, error) {
	var res Response

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	updateFields["updated_at"] = time.Now().In(loc)
	updated_at := updateFields["updated_at"]

	con := db.CreateCon()

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

	sqlStatement := "UPDATE user " + setStatement + " WHERE user_id = ?"
	values = append(values, userID)

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

// DeleteUser deletes a user with the provided ID
func DeleteUser(userID int) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "DELETE FROM user WHERE user_id = ?"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(userID)

	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected":    rowsAffected,
		"deleted_user_id": userID,
	}

	return res, err
}
