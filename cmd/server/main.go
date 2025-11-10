package main

import (
	"log"

	"booking.com/internal/config"
	"booking.com/internal/db/postgresql/dao"
	"booking.com/internal/server"
)

func main() {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		log.Println("error in loading app configuration, error: ", err)
		return
	}

	db, err := dao.Connect(cfg.PostgresqlDb)
	if err != nil {
		log.Println("error in connecting db, error: ", err)
		return
	}
	dao.SetDefault(db)
	if err := server.StartHttpTlsServer(cfg.HttpServer); err != nil {
		log.Printf("server failed, error: %v", err)
	}
}
