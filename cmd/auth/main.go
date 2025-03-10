package main

import (
	"log"

	auth "github.com/fatykhovar/jwtAuth/internal/auth"
	"github.com/fatykhovar/jwtAuth/internal/config"
	"github.com/fatykhovar/jwtAuth/internal/service"
	"github.com/fatykhovar/jwtAuth/internal/storage"
	"github.com/fatykhovar/jwtAuth/internal/storage/postgres"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustLoad()

	// подключение к бд
	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	storage := storage.NewStorage(db)
	service := service.NewService(storage)
	// запуск сервера
	server := auth.NewAPIServer(":3000")
	server.Run()

	i := 0
	if err := server.InsertTestData(auth.IP_address[i]); err != nil {
		log.Fatal(err)
	}
}
