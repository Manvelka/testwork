package gender_test

import (
	"context"
	"testing"

	"github.com/Manvelka/testwork/pkg/enrich/gender"
)

func TestDefaultApi(t *testing.T) {
	if _, err := gender.DefaultGenderApi.EnrichGender(context.Background(), "Alex"); err != nil {
		t.Error(err)
	}
}
