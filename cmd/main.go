package main

import (
	"fmt"
	"subscription/internal"
)

func main() {

	subsMap := internal.NewSubStore()
	handlers := internal.NewHTTPHandlers(&subsMap)
	server := internal.NewHTTPServer(handlers)
	if err := server.StartServer(); err != nil {
		fmt.Println(err)
		return
	}

}
