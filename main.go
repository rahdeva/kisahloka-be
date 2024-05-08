package main

import (
	"kisahloka_be/db"
	"kisahloka_be/routes"
)

func main() {
	db.DBInit()

	e := routes.Init()

	e.Logger.Fatal(e.Start(":4000"))
}
