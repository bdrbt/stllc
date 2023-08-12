package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bdrbt/stllc/pkg/dto"
)

func (svc *Service) syncData(ctx context.Context) error {
	log.Print("starting update")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, svc.feedURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request:%w", err)
	}

	req.Header.Set("Accept", "application/xml")

	client := &http.Client{}

	rsp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error requsting remote url:%w", err)
	}

	defer rsp.Body.Close()

	// yep, i know it's huge
	raw, err := io.ReadAll(rsp.Body)
	if err != nil {
		return fmt.Errorf("error reading response stream:%w", err)
	}

	responseData := dto.SdnResponse{}

	err = xml.Unmarshal(raw, &responseData)
	if err != nil {
		return fmt.Errorf("error deserialising remote response:%w", err)
	}

	for _, rec := range responseData.SdnEntries {
		err := svc.repository.Upsert(ctx, rec.Domain())
		if err != nil {
			return fmt.Errorf("cannot upsert record:%w", err)
		}

		log.Printf("Entity:%s", rec.Pretty())
	}

	// log update
	log.Printf("received %d records", len(responseData.SdnEntries))

	return nil
}
