package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"

	"subs-service/internal/config"
	myhttp "subs-service/internal/http"
	"subs-service/internal/repository"
	"subs-service/internal/service"
	"subs-service/logger"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Логгер
	lg := logger.New()

	// Формируем строку подключения к БД из env
	dsn := "host=" + cfg.DBHost +
		" port=" + cfg.DBPort +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" sslmode=" + cfg.DBSSLMode

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Репозиторий
	repo := repository.NewPostgresSubscriptionRepository(db)

	// Сервис
	srv := service.NewSubscriptionService(repo, lg)

	// HTTP handler
	h := myhttp.NewHandler(srv)
	h.RegisterRoutes()

	// Запуск сервера
	log.Println("Starting server on port", cfg.AppPort)
	http.ListenAndServe(":"+cfg.AppPort, nil)
}
