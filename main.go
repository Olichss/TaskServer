package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func main() {
	server, err := NewTaskServer()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer server.DB.Close()

	log.Info("Starting task server...")

	router := mux.NewRouter()
	router.HandleFunc("/task/", server.CreateTaskHandler).Methods("POST")
	router.HandleFunc("/task/{id:[0-9]+}/", server.GetTaskHandler).Methods("GET")
	router.HandleFunc("/task/{id:[0-9]+}/", server.DeleteTaskHandler).Methods("DELETE")
	router.HandleFunc("/task/", server.GetAllTaskHandler).Methods("GET")
	router.HandleFunc("/task/", server.DeleteAllTaskHandler).Methods("DELETE")

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	err = http.ListenAndServe(":8000", handler)
	if err != nil {
		fmt.Print(err)
	}
}
