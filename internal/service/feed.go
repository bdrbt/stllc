package service

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bdrbt/stllc/pkg/dto"
	"github.com/google/uuid"
)

func (svc *Service) syncData() error {
	updID := uuid.New()

	log.Print("starting update")
	req, err := http.NewRequest(http.MethodGet, svc.feedURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/xml")

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}

	// yep, i know it's huge
	// TODO optimise token reading
	raw, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	responseData := dto.SdnResponse{}
	xml.Unmarshal(raw, &responseData)
	for _, rec := range responseData.SdnEntries {
		err := svc.repository.Upsert(rec.Domain(), updID.String())
		if err != nil {
			return fmt.Errorf("cannot upsert record:%v", err)
		}
		log.Printf("Entity:%s", rec.Pretty())
	}

	// log update
	log.Printf("received %d records", len(responseData.SdnEntries))

	return nil
}
