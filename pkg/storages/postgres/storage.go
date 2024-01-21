package storage

import (
	"context"
	"database/sql"

	"github.com/Manvelka/testwork/pkg/person"
)

type Storage struct {
	DB *sql.DB
}

const migration = `
create table if not exists persons (
	person_id  serial primary key,
	name       text    not null,
	surname    text    not null,
	patronymic text,
	age        integer,
	gender     text,
	nation     text
);
`

func (s *Storage) Migration(ctx context.Context) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if _, err := tx.Exec(migration); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Storage) Get(ctx context.Context) ([]person.Person, error) {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, "select person_id, name, surname, patronymic, age, gender, nation from persons;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	persons := []person.Person{}
	for rows.Next() {
		var (
			p                          person.Person
			patronymic, gender, nation sql.NullString
			age                        sql.NullInt64
		)
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname, &patronymic, &age, &gender, &nation); err != nil {
			return nil, err
		}
		if patronymic.Valid {
			p.Patronymic = patronymic.String
		}
		if gender.Valid {
			p.Gender = gender.String
		}
		if nation.Valid {
			p.Nation = nation.String
		}
		if age.Valid {
			p.Age = int(age.Int64)
		}
		persons = append(persons, p)
	}
	return persons, tx.Commit()
}

func (s *Storage) Post(ctx context.Context, p person.Person) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(
		ctx,
		"insert into persons (name, surname, patronymic, age, gender, nation) values ($1, $2, $3, $4, $5, $6);",
		p.Name,
		p.Surname,
		p.Patronymic,
		p.Age,
		p.Gender,
		p.Nation,
	); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Storage) Put(ctx context.Context, p person.Person) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(
		ctx,
		"update persons set name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nation = $6 where person_id = $7;",
		p.Name,
		p.Surname,
		p.Patronymic,
		p.Age,
		p.Gender,
		p.Nation,
		p.ID,
	); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Storage) Delete(ctx context.Context, personID int) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "delete from persons where person_id = $1;", personID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
