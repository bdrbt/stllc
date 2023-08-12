package domain

import "time"

type SDNRecord struct {
	UID       int64  `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
