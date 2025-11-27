package main

import (
	"log"
	"net/http"
	"time"
	"trinity/internal/config"
	"trinity/internal/handler"
	"trinity/internal/middleware"
	"trinity/internal/repository_postgres"
	"trinity/internal/repository_redis"
	"trinity/pkg/jwt"
	"trinity/pkg/smtp"
)

func main() {
	mux := http.NewServeMux()
	authMux := http.NewServeMux()
	softAuthMux := http.NewServeMux()
	adminMux := http.NewServeMux()
	config.Init()
	err := jwt.Init()
	if err != nil {
		log.Fatal(err)
	}
	err = smtp.InitSMTP()
	if err != nil {
		log.Fatal(err)
	}
	db := &repository_postgres.PostgresDB{}
	err = db.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}
	rdb := repository_redis.RedisDB{}
	err = rdb.ConnectToRedis()
	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(db)

	mux.HandleFunc("/", h.MainPage)
	mux.HandleFunc("/static/", h.Static)
	mux.HandleFunc("POST /decode/{algorithm}", h.Decode)
	mux.HandleFunc("POST /login", h.LoginUser)
	mux.HandleFunc("POST /register", h.RegisterUser)
	mux.Handle("/api/", middleware.Auth(authMux))
	mux.Handle("/encode/", middleware.SoftAuth(softAuthMux))
	mux.Handle("/admin/", middleware.Auth(middleware.AdminAuth(adminMux)))
	authMux.HandleFunc("GET /api/history/{algorithm}", h.GetHistory)
	authMux.HandleFunc("DELETE /api/delete-user", h.DeleteUser)
	authMux.HandleFunc("DELETE /api/delete-user/{id}", h.DeleteUser)
	softAuthMux.HandleFunc("POST /encode/{algorithm}", h.Encode)
	adminMux.HandleFunc("GET /admin/history/{algorithm}/{login}", h.GetHistoryByLogin)
	adminMux.HandleFunc("GET /admin/history/", h.GetHistoryPage)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Println("Service working on port http://localhost:8080/")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
