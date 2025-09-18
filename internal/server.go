// server.go роутер для реализации RESTfullAPI
package internal

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Структура для работы с нашими хендлерами
type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(httpHandlers *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandlers,
	}
}

func (s *HTTPServer) StartServer() error {
	r := mux.NewRouter()
	r.HandleFunc("/subscriptions", s.httpHandlers.HandleSubscribe).Methods("POST")
	r.HandleFunc("/subscriptions/", s.httpHandlers.HandleGetAllInfoSubscribe).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", s.httpHandlers.HandleGetInfoSubscribe).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", s.httpHandlers.HandleDeleteSubscribe).Methods("DELETE")
	r.HandleFunc("/subscriptions/{id}", s.httpHandlers.HandleUpdateSubscribe).Methods("PUT")
	// r.Path("/subscriptions").Methods("POST").HandlerFunc(s.httpHandlers.HandleSubscribe)
	// r.Path("/subscriptions/{id}").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetInfoSubscribe)
	// r.Path("/subscriptions/{id}").Methods("DELETE").HandlerFunc(s.httpHandlers.HandleDeleteSubscribe)

	return http.ListenAndServe(":8080", r)
}
