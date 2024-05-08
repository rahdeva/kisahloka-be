package db

import (
	"database/sql"
	"fmt"
	"kisahloka_be/config"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func DBInit() {
	DBconf := config.GetDBConfig()

	// username:password@protocol(address)/dbname?param=value
	DBconnectionString := DBconf.DB_USERNAME + ":" + DBconf.DB_PASSWORD + "@tcp(" + DBconf.DB_HOST + ":" + DBconf.DB_PORT + ")/" + DBconf.DB_NAME

	db, err = sql.Open("mysql", DBconnectionString+"?parseTime=true")

	if err != nil {
		panic("Connection Error")
	}

	err = db.Ping()

	if err != nil {
		fmt.Println(DBconnectionString)
		panic("DSN Error")
	}
}

func CreateCon() *sql.DB {
	return db
}
