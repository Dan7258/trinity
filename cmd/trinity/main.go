package main

import (
	"log"
	"net/http"
	"time"
	"trinity/internal/config"
	"trinity/internal/handler"
	"trinity/internal/middleware"
	"trinity/internal/repository"
	"trinity/pkg/jwt"
)

func main() {
	mux := http.NewServeMux()
	authMux := http.NewServeMux()
	softAuthMux := http.NewServeMux()
	config.Init()
	err := jwt.Init()
	if err != nil {
		log.Fatal(err)
	}
	db := &repository.PostgresDB{}
	err = db.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(db)

	mux.HandleFunc("/", h.MainPage)
	mux.HandleFunc("/static/", h.Static)
	mux.HandleFunc("POST /decode/{algorithm}", h.Decode)
	mux.HandleFunc("POST /login", h.LoginUser)
	mux.HandleFunc("POST /register", h.RegisterUser)
	mux.Handle("/history/", middleware.Auth(authMux))
	mux.Handle("/encode/", middleware.SoftAuth(softAuthMux))
	authMux.HandleFunc("GET /history/{algorithm}", h.GetHistory)
	softAuthMux.HandleFunc("POST /encode/{algorithm}", h.Encode)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
