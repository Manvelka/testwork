package age

import (
	"context"
	"encoding/json"
	"net/http"
)

var DefaultAgeApi = &AgeApi{Host: "https://api.agify.io"}

type enrichAge struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type AgeApi struct {
	Host string
}

func (a *AgeApi) EnrichAge(ctx context.Context, name string) (age int, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.Host+"?name="+name, nil)
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var o enrichAge
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return 0, err
	}
	return o.Age, nil
}
