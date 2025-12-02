# Subscriptions Service (Effective Mobile Test Task)

REST-сервис для хранения данных о пользовательских подписках и расчёта
их стоимости за период.\
Решение выполнено согласно тестовому заданию Effective Mobile.

## Функциональность

-   CRUD операции над подписками
-   Вывод списка подписок с фильтрами
-   Расчёт суммарной стоимости подписок за период
-   Миграции PostgreSQL
-   Логирование
-   Конфигурация через `.env`
-   Swagger документация
-   Полная Docker-сборка (`docker-compose up`)

## Технологии

-   Go 1.24+
-   PostgreSQL 15
-   net/http
-   Docker, Docker Compose
-   OpenAPI 3 (Swagger)

## Структура проекта

    subs-service/
    ├── cmd/
    │   └── sub-service/
    │       └── main.go
    ├── internal/
    │   ├── config/
    │   │   └── config.go
    │   ├── db/
    │   │   └── migrations/
    │   │            ├── 001_init.up.sql
    │   │            └── 001_init.down.sql
    │   ├── domain/
    │   │   └── subscription.go
    │   ├── repository/
    │   │   ├── subs_repo.go
    │   │   └── subs_repo_interface.go
    │   ├── service/
    │   │   └── subs_service.go
    │   └── http/
    │       ├── handlers.go
    │       ├── router.go
    │       └── swagger.go
    ├── docs/
    │   └── swagger.yaml
    ├── logger/
    │   └── logger.go
    ├── Dockerfile
    ├── docker-compose.yml
    └── .env

## Переменные окружения

    APP_PORT=8080
    DB_HOST=db
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=subs
    DB_SSLMODE=disable

## Запуск

    docker-compose up -d --build

Сервис будет доступен по адресу http://localhost:8080

## Swagger

Документация доступна по адресу:

    http://localhost:8080/docs/ - документация с UI
    http://localhost:8080/swagger - код
