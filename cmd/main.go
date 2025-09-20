package main

import "subscription/internal/database"

func main() {

	db := database.NewDBConnection()
	defer db.Close()
}
