package dao

import (
	"fmt"
	"log"

	"booking.com/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBS struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Schema   string
}

func Connect(cfg config.PostgreSQL) (*gorm.DB, error) {
	d := Init(cfg)
	db, err := gorm.Open(postgres.Open(d.connectionString()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println("error in Open ", err)
		return nil, err
	}

	return db, nil
}

func (d DBS) connectionString() string {
	return fmt.Sprintf("host=%v port=%v dbname=%v user=%v search_path=%s password=%v sslmode=disable", d.Host, d.Port, d.Name, d.User, d.Schema, d.Password)
}

func Init(cfg config.PostgreSQL) DBS {
	return DBS{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Name:     cfg.Name,
		User:     cfg.UserName,
		Password: cfg.Password,
		Schema:   cfg.Schema,
	}
}
