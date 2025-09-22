package server

import (
	"fmt"
	"net/http"
	"os"
	"subscription/internal/handlers"

	"github.com/gorilla/mux"
)

// Структура для работы с нашими хендлерами
type HTTPServer struct {
	httpHandlers *handlers.HTTPHandlers
}

func NewHTTPServer(httpHandlers *handlers.HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandlers,
	}
}

func (s *HTTPServer) StartServer() error {

	port := os.Getenv("SERVER_PORT")

	r := mux.NewRouter()
	r.HandleFunc("/subscriptions", s.httpHandlers.HandleSubscribe).Methods("POST")
	r.HandleFunc("/subscriptions", s.httpHandlers.HandleGetAllInfoSubscribe).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", s.httpHandlers.HandleGetInfoSubscribe).Methods("GET")
	r.HandleFunc("/subscriptions/sum", s.httpHandlers.HandleSumInfo).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", s.httpHandlers.HandleDeleteSubscribe).Methods("DELETE")
	r.HandleFunc("/subscriptions/{id}", s.httpHandlers.HandleUpdateSubscribe).Methods("PUT")
	// r.Path("/subscriptions").Methods("POST").HandlerFunc(s.httpHandlers.HandleSubscribe)
	// r.Path("/subscriptions/{id}").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetInfoSubscribe)
	// r.Path("/subscriptions/{id}").Methods("DELETE").HandlerFunc(s.httpHandlers.HandleDeleteSubscribe)
	fmt.Println("Start Server")
	fmt.Println("port", port)
	return http.ListenAndServe(port, r)
}
