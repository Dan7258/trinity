package main

import (
	"fmt"
	"log"
	"time"
	"trinity/pkg/big_int_generate"
)

func main() {
	ch := make(chan bool)
	go spinner(ch)
	n := big_int_generate.GetRandomPrime(4096)
	ch <- true
	close(ch)
	log.Println(n.String())

	//mux := http.NewServeMux()
	//authMux := http.NewServeMux()
	//softAuthMux := http.NewServeMux()
	//adminMux := http.NewServeMux()
	//config.Init()
	//err := jwt.Init()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = smtp.InitSMTP()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//db := &repository_postgres.PostgresDB{}
	//err = db.ConnectToDatabase()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//rdb := repository_redis.RedisDB{}
	//err = rdb.ConnectToRedis()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//bot := telegram.NewBot(db, rdb)
	//
	//err = bot.ConnectBot()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//go bot.HandleUpdates()
	//
	//h := handler.NewHandler(db, rdb)
	//
	//mux.HandleFunc("/", h.MainPage)
	//mux.HandleFunc("/static/", h.Static)
	//mux.HandleFunc("POST /decode/{algorithm}", h.Decode)
	//mux.HandleFunc("POST /login", h.LoginUser)
	//mux.HandleFunc("POST /register", h.RegisterUser)
	//mux.HandleFunc("POST /secure/two-fa/{email}", h.TwoFA)
	//mux.Handle("/api/", middleware.Auth(authMux))
	//mux.Handle("/encode/", middleware.SoftAuth(softAuthMux))
	//mux.Handle("/admin/", middleware.Auth(middleware.AdminAuth(adminMux)))
	//authMux.HandleFunc("GET /api/history/{algorithm}", h.GetHistory)
	//authMux.HandleFunc("DELETE /api/delete-user", h.DeleteUser)
	//authMux.HandleFunc("DELETE /api/delete-user/{id}", h.DeleteUser)
	//softAuthMux.HandleFunc("POST /encode/{algorithm}", h.Encode)
	//adminMux.HandleFunc("GET /admin/history/{algorithm}/{login}", h.GetHistoryByLogin)
	//adminMux.HandleFunc("GET /admin/history/", h.GetHistoryPage)
	//server := &http.Server{
	//	Addr:         ":8080",
	//	Handler:      mux,
	//	ReadTimeout:  60 * time.Second,
	//	WriteTimeout: 60 * time.Second,
	//	IdleTimeout:  120 * time.Second,
	//}
	//log.Println("Service working on port http://localhost:8080/")
	//err = server.ListenAndServe()
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func spinner(done chan bool) {
	chars := []rune{'|', '/', '-', '\\'}
	start := time.Now()
	for {
		select {
		case <-done:
			fmt.Println("\nчисла сгенерировались")
			return
		default:
			for _, r := range chars {
				now := time.Since(start).Seconds()
				fmt.Printf("\r%c прошло %v сек.", r, int(now))
				time.Sleep(100 * time.Millisecond)
			}

		}
	}
}
