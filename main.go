package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/Manvelka/testwork/app"
	"github.com/Manvelka/testwork/pkg/enrich/age"
	"github.com/Manvelka/testwork/pkg/enrich/gender"
	"github.com/Manvelka/testwork/pkg/enrich/nation"
	storage "github.com/Manvelka/testwork/pkg/storages/postgres"

	_ "github.com/lib/pq"
)

func config() (dbInit string, err error) {
	if err := godotenv.Load(".env"); err != nil {
		return "", err
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	db := os.Getenv("DBNAME")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, db), nil
}

func main() {
	postgresInit, err := config()
	if err != nil {
		log.Fatalf("Ошибка чтения конфигурационного файла: %v", err)
	}

	db, err := sql.Open("postgres", postgresInit)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("ошибка подключения: %v", err)
	}

	app := &app.App{
		AgeService:     age.DefaultAgeApi,
		GenderService:  gender.DefaultGenderApi,
		NationService:  nation.DefaultNationApi,
		StorageService: &storage.Storage{DB: db},
		InfoLogger:     *log.New(os.Stdout, "INFO: ", log.Ltime),
		ErrLogger:      *log.New(os.Stderr, "ERROR: ", log.Ltime),
	}

	app.StorageService.Migration(context.Background())

	if err := http.ListenAndServe(":5000", app); err != nil {
		app.ErrLogger.Printf("server error: %v", err)
	}
}
