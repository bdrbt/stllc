package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi"

	"github.com/bdrbt/stllc/internal/config"
	"github.com/bdrbt/stllc/internal/repository"
	"github.com/bdrbt/stllc/pkg/dto"
)

const (
	Starting = iota
	Ready
	Busy
)

type Service struct {
	state      int                    // Current service state, on of Starting, Ready, Busy
	lock       sync.RWMutex           // mutex to handle state flag switching
	router     *chi.Mux               // router
	feedURL    string                 // feed URL
	repository *repository.Repository // storage
}

// New - creating new Service instance
func New(cfg *config.Config) (*Service, error) {
	repo, err := repository.New(cfg.PgURL())
	if err != nil {
		return nil, fmt.Errorf("cannot create repository:%w", err)
	}

	svc := &Service{
		feedURL:    "https://www.treasury.gov/ofac/downloads/sdn.xml",
		state:      Starting,
		router:     chi.NewRouter(),
		repository: repo,
	}

	// setup soutes
	svc.router.Get("/state", svc.State)
	svc.router.Get("/update", svc.Update)
	svc.router.Get("/get_names", svc.GetNames)

	return svc, nil
}

// Run - start processing http request on congfigured port,
// see internal/conf package for details.
func (svc *Service) Run() error {
	svc.lock.Lock()
	svc.state = Ready
	svc.lock.Unlock()

	go http.ListenAndServe(":8080", svc.router)
	return nil
}

// JSONResponse - serialize provided data into JSON and write them into ResponseWriter
func (svc *Service) JSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json, err := json.Marshal(data)
	if err != nil {
		log.Printf("error writing response:%v", err)
	}
	w.Write(json)
}

func (svc *Service) UpdateSuccessResponse(w http.ResponseWriter) {
	svc.JSONResponse(w, dto.UpdateResponse{
		Result: true,
		Code:   200,
		Info:   "",
	})
}

func (svc *Service) UpdateFailResponse(w http.ResponseWriter) {
	svc.JSONResponse(w, dto.UpdateResponse{
		Result: false,
		Code:   503,
		Info:   "service unavailable",
	})
	return
}
