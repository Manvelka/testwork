package gender

import (
	"context"
	"encoding/json"
	"net/http"
)

var DefaultGenderApi = &GenderApi{Host: "https://api.genderize.io"}

type enrichGender struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float32 `json:"probability"`
}

type GenderApi struct {
	Host string
}

func (a *GenderApi) EnrichGender(ctx context.Context, name string) (gender string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.Host+"?name="+name, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var o enrichGender
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return "", err
	}
	return o.Gender, nil
}
