package nation_test

import (
	"context"
	"testing"

	"github.com/Manvelka/testwork/pkg/enrich/nation"
)

func TestDefaultApi(t *testing.T) {
	if _, err := nation.DefaultNationApi.EnrichNation(context.Background(), "Alex"); err != nil {
		t.Error(err)
	}
}
