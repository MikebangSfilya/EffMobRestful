package main

import (
	"log"
	"subscription/internal/api/handlers"
	"subscription/internal/api/server"
	"subscription/internal/database"
	"subscription/internal/repository"
	"subscription/internal/service"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db := database.NewDBConnection()
	defer db.Close()

	repo := repository.NewPgxRepository(db)

	serv := service.NewService(repo)

	h := handlers.NewHTTPHandlers(serv)

	srv := server.NewHTTPServer(h)
	if err := srv.StartServer(); err != nil {
		log.Printf("Internal Server problem %v", err)
		return
	}

}
