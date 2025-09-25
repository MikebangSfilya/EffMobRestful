package main

import (
	"log"
	"subscription/internal/database"
	"subscription/internal/handlers"
	"subscription/internal/model"
	"subscription/internal/server"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db := database.NewDBConnection()
	defer db.Close()
	store := model.NewSubStore(db)
	h := handlers.NewHTTPHandlers(store)
	srv := server.NewHTTPServer(h)
	if err := srv.StartServer(); err != nil {
		log.Printf("Internal Server problem %v", err)
		return
	}

}
