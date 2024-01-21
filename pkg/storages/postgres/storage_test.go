package storage_test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/Manvelka/testwork/pkg/person"
	storage "github.com/Manvelka/testwork/pkg/storages/postgres"
)

const testingPostgres = "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"

func TestMigration(t *testing.T) {
	db, err := sql.Open("postgres", testingPostgres)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	s := storage.Storage{DB: db}
	if err := s.Migration(context.Background()); err != nil {
		t.Error(err)
	}
}

func TestPost(t *testing.T) {
	db, err := sql.Open("postgres", testingPostgres)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	s := storage.Storage{DB: db}
	if err := s.Post(context.Background(), person.Person{
		Name:       "Alex",
		Surname:    "Tsiporin",
		Patronymic: "Sergeevich",
		Age:        40,
		Gender:     "male",
		Nation:     "RU",
	}); err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	db, err := sql.Open("postgres", testingPostgres)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	s := storage.Storage{DB: db}
	pp, err := s.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(pp) == 0 {
		t.Fatal("нет данных в таблице")
	}
	if pp[0].Name != "Alex" {
		t.Error(pp[0].Name)
	}
}

func TestPut(t *testing.T) {
	db, err := sql.Open("postgres", testingPostgres)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	s := storage.Storage{DB: db}
	if err := s.Put(context.Background(), person.Person{
		ID:         1,
		Name:       "Александр",
		Surname:    "Ципорин",
		Patronymic: "Сергеевич",
		Age:        40,
		Gender:     "male",
		Nation:     "RU",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	db, err := sql.Open("postgres", testingPostgres)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	s := storage.Storage{DB: db}
	if err := s.Delete(context.Background(), 1); err != nil {
		t.Error(err)
	}
}
