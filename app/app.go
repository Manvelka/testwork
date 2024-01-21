package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Manvelka/testwork/pkg/person"
)

type App struct {
	AgeService     AgeEnricher
	GenderService  GenderEnricher
	NationService  NationEnricher
	StorageService Storager
	InfoLogger     log.Logger
	ErrLogger      log.Logger
}

func (a *App) enrich(ctx context.Context, p *person.Person) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	m := &sync.Mutex{}
	go func() {
		defer wg.Done()
		age, err := a.AgeService.EnrichAge(ctx, p.Name)
		if err != nil {
			a.ErrLogger.Printf("error with enrich %s age: %v", p.Name, err)
		} else {
			a.InfoLogger.Printf("Обогащение имени %s возрастом: %d", p.Name, age)
			m.Lock()
			p.Age = age
			m.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		gender, err := a.GenderService.EnrichGender(ctx, p.Name)
		if err != nil {
			a.ErrLogger.Printf("error with enrich %s gender: %v", p.Name, err)
		} else {
			a.InfoLogger.Printf("Обогащение имени %s полом: %s", p.Name, gender)
			m.Lock()
			p.Gender = gender
			m.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		nation, err := a.NationService.EnrichNation(ctx, p.Name)
		if err != nil {
			a.ErrLogger.Printf("error with enrich %s nation: %v", p.Name, err)
		} else {
			a.InfoLogger.Printf("Обогащение имени %s национальностью: %s", p.Name, nation)
			m.Lock()
			p.Nation = nation
			m.Unlock()
		}
	}()
	wg.Wait()
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.Get(w, r)
	case http.MethodPost:
		a.Post(w, r)
	case http.MethodPut:
		a.Put(w, r)
	case http.MethodDelete:
		a.Delete(w, r)
	default:
		http.Error(w, "not valid method", http.StatusBadRequest)
	}
}

func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	pp, err := a.StorageService.Get(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("error with select objects: %v", err), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(pp); err != nil {
		http.Error(w, fmt.Sprintf("error with encode: %v", err), http.StatusInternalServerError)
		return
	}
	a.InfoLogger.Printf("GET %s", r.RemoteAddr)
}

func (a *App) Post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var p person.Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, fmt.Sprintf("error with decode struct: %v", err), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	a.enrich(ctx, &p)
	if err := a.StorageService.Post(ctx, p); err != nil {
		http.Error(w, fmt.Sprintf("error with create object: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	a.InfoLogger.Printf("POST %s: %v", r.RemoteAddr, p)
}

func (a *App) Put(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var p person.Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, fmt.Sprintf("error with decode struct: %v", err), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	if err := a.StorageService.Put(ctx, p); err != nil {
		http.Error(w, fmt.Sprintf("error with update object: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	a.InfoLogger.Printf("PUT %s: %v", r.RemoteAddr, p)
}

func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var p person.Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, fmt.Sprintf("error with decode struct: %v", err), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	if err := a.StorageService.Delete(ctx, p.ID); err != nil {
		http.Error(w, fmt.Sprintf("error with delete object: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	a.InfoLogger.Printf("DELETE %s: %v", r.RemoteAddr, p)
}
