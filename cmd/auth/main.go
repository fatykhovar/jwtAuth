package main

import (
	// "database/sql"
	"fmt"
	"log"

	auth "github.com/fatykhovar/jwtAuth/internal/auth"
	"github.com/fatykhovar/jwtAuth/internal/storage"
	_ "github.com/lib/pq"
)

func main() {
	// подключение к бд
	store, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	// создание таблицы токенов
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", store)

	// запуск сервера
	server := auth.NewAPIServer(":3000", store)
	server.Run()

	i := 0
	if err := server.InsertTestData(auth.IP_address[i]); err != nil {
		log.Fatal(err)
	}
}
