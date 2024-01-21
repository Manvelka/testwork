package app

import (
	"context"

	"github.com/Manvelka/testwork/pkg/person"
)

type AgeEnricher interface {
	EnrichAge(ctx context.Context, name string) (age int, err error)
}

type GenderEnricher interface {
	EnrichGender(ctx context.Context, name string) (gender string, err error)
}

type NationEnricher interface {
	EnrichNation(ctx context.Context, name string) (nation string, err error)
}

type Storager interface {
	Migration(ctx context.Context) error
	Get(ctx context.Context) ([]person.Person, error)
	Post(ctx context.Context, p person.Person) error
	Put(ctx context.Context, p person.Person) error
	Delete(ctx context.Context, personID int) error
}
