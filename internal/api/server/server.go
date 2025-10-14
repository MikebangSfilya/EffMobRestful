package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Структура для работы с нашими хендлерами
type HTTPServer struct {
	httpHandlers HTTPRepository
}

type HTTPRepository interface {
	HandleSubscribe(w http.ResponseWriter, r *http.Request)
	HandleGetInfoSubscribe(w http.ResponseWriter, r *http.Request)
	HandleGetAllInfoSubscribe(w http.ResponseWriter, r *http.Request)
	HandleDeleteSubscribe(w http.ResponseWriter, r *http.Request)
	HandleUpdateSubscribe(w http.ResponseWriter, r *http.Request)
	HandleSumInfo(w http.ResponseWriter, r *http.Request)
}

func NewHTTPServer(httpHandlers HTTPRepository) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandlers,
	}
}

func (s *HTTPServer) StartServer() error {

	port := os.Getenv("SERVER_PORT")
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	r := chi.NewRouter()
	r.Post("/subscriptions", s.httpHandlers.HandleSubscribe)
	r.Get("/subscriptions", s.httpHandlers.HandleGetAllInfoSubscribe)
	r.Get("/subscriptions/{id}", s.httpHandlers.HandleGetInfoSubscribe)
	r.Get("/subscriptions/sum", s.httpHandlers.HandleSumInfo)
	r.Delete("/subscriptions/{id}", s.httpHandlers.HandleDeleteSubscribe)
	r.Put("/subscriptions/{id}", s.httpHandlers.HandleUpdateSubscribe)
	fmt.Println("Start Server")
	fmt.Println("port", port)
	return http.ListenAndServe(port, r)
}
