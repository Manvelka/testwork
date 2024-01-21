package age_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Manvelka/testwork/pkg/enrich/age"
)

func mock(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !r.Form.Has("name") {
		http.Error(w, "отсутствует параметр", http.StatusBadRequest)
		return
	}
	if r.Form.Get("name") != "Alex" {
		http.Error(w, "переданное имя не соответствует условию", http.StatusBadRequest)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "метод не соответствует условию", http.StatusBadRequest)
		return
	}
	w.Write([]byte(`{"count":1,"name":"Alex","age":40}`))
}

func TestEnrichAgeMoch(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(mock))
	defer s.Close()
	m := &age.AgeApi{Host: s.URL}
	age, err := m.EnrichAge(context.TODO(), "Alex")
	if err != nil {
		t.Error(err)
	}
	if age != 40 {
		t.Error(age)
	}
}

func TestDefaultApi(t *testing.T) {
	if _, err := age.DefaultAgeApi.EnrichAge(context.Background(), "Alex"); err != nil {
		t.Error(err)
	}
}
