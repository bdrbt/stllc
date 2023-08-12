package service

import (
	"log"
	"net/http"

	"github.com/bdrbt/stllc/internal/domain"
)

// State handler.
//
//nolint:revive
func (svc *Service) State(w http.ResponseWriter, r *http.Request) {
	svc.lock.RLock()
	defer svc.lock.RUnlock()

	switch svc.state {
	case Starting:
		c, err := w.Write([]byte("{\"result\": false, \"info\": \"service is starting\"}"))
		if err != nil || c == 0 {
			log.Printf("error writing response:%v", err)

			return
		}
	case Busy:
		c, err := w.Write([]byte("{\"result\": false, \"info\": \"data is updating\"}"))
		if err != nil || c == 0 {
			log.Printf("error writing response:%v", err)

			return
		}

	default:
		c, err := w.Write([]byte("{\"result\": true, \"info\": \"ok\"}"))
		if err != nil || c == 0 {
			log.Printf("error writing response:%v", err)

			return
		}
	}
}

// Update handler.
func (svc *Service) Update(w http.ResponseWriter, r *http.Request) {
	// Check if service is busy or starting.
	svc.lock.Lock()
	if svc.state != Ready {
		// releas lock and finish.
		svc.lock.Unlock()
		svc.UpdateFailResponse(w)

		return
	}

	// set it's into Busy state.
	svc.state = Busy
	svc.lock.Unlock()

	// ensure we change service state at the end.
	defer func() {
		svc.lock.Lock()
		svc.state = Ready
		svc.lock.Unlock()
	}()

	err := svc.syncData(r.Context())
	if err != nil {
		log.Printf("update error:%v", err)
		svc.UpdateFailResponse(w)
	} else {
		svc.UpdateSuccessResponse(w)
	}
}

// GenNames handler.
func (svc *Service) GetNames(w http.ResponseWriter, r *http.Request) {
	var (
		recs []domain.SDNRecord
		err  error
	)

	name := r.URL.Query().Get("name")
	queryType := r.URL.Query().Get("type")

	if name == "" {
		log.Printf("name param not provided")

		return
	}

	if queryType == "strong" {
		log.Printf("strong quereyng:\"%s\"", name)

		recs, err = svc.repository.QueryByName(r.Context(), name)
		if err != nil {
			log.Printf("error searching by name")

			return
		}
	} else {
		log.Printf("weak quereyng:\"%s\"", name)
		recs, err = svc.repository.QueryByPattern(r.Context(), name)

		if err != nil {
			log.Printf("error searching by pattern")

			return
		}
	}

	svc.JSONResponse(w, recs)
}
