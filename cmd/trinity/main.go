package main

import (
	"log"
	"net/http"
	"time"
	"trinity/internal/config"
	"trinity/internal/handler"
	"trinity/internal/model"
)

func main() {
	mux := http.NewServeMux()
	config.Init()
	db := &model.PostgresDB{}
	err := db.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(db)

	mux.HandleFunc("/", h.MainPage)
	mux.HandleFunc("/static/", h.Static)
	mux.HandleFunc("POST /encode/{algorithm}", h.Encode)
	mux.HandleFunc("POST /decode/{algorithm}", h.Decode)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
