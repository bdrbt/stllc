package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"github.com/bdrbt/stllc/internal/domain"
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

	records := make([]domain.SDNRecord, 0)

	decoder := xml.NewDecoder(rsp.Body)
	for {
		t, _ := decoder.Token()

		// break on EOF.
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "sdnEntry" {
				var rec dto.SdnEntry

				decoder.DecodeElement(&rec, &se)
				if rec.SdnType == "Individual" {
					records = append(records, rec.Domain())
				}
			}
		default:
		}
	}

	rsp.Body.Close()

	for _, rec := range records {
		err := svc.repository.Upsert(ctx, rec)
		if err != nil {
			return fmt.Errorf("cannot upsert record:%w", err)
		}
	}

	log.Printf("received %d individual records", len(records))

	return nil
}
