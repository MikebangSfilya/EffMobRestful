package main

import (
	"subscription/internal"
)

func main() {

	db := internal.NewDBConnection()
	defer db.Close()
}
