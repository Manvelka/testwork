package nation

import (
	"context"
	"encoding/json"
	"net/http"
)

var DefaultNationApi = &NationApi{Host: "https://api.nationalize.io"}

type country struct {
	CountryID   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type enrichNation struct {
	Count   int       `json:"count"`
	Name    string    `json:"name"`
	Country []country `json:"country"`
}

func (e *enrichNation) counryWithMaxProbability() string {
	return e.Country[0].CountryID
}

type NationApi struct {
	Host string
}

func (a *NationApi) EnrichNation(ctx context.Context, name string) (nation string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.Host+"?name="+name, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var o enrichNation
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return "", err
	}
	return o.counryWithMaxProbability(), nil
}
