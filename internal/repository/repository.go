package repository

import (
	"context"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bdrbt/stllc/internal/domain"
)

const (
	upsertSQL = `
		INSERT INTO sdn_records (id, first_name, last_name, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		ON CONFLICT (id)
		DO UPDATE SET
		first_name = $2, last_name = $3, updated_at = now()
	`

	queryStrictSQL = `
		SELECT id,first_name,last_name,created_at,updated_at FROM sdn_records
		WHERE lower(first_name) = $1 OR LOWER(last_name) = $1
	`

	queryWeakSQL = `
		SELECT id,first_name,last_name,created_at,updated_at FROM sdn_records
		WHERE first_name ILIKE ANY ($1) or last_name ilike ANY($1)
	`
)

type Repository struct {
	db *pgxpool.Pool
}

func New(DBURL string) (*Repository, error) {
	repo := &Repository{}

	poolConfig, err := pgxpool.ParseConfig(DBURL)
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_URL:", err)
	}

	repo.db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}

	return repo, nil
}

// Upsert - insert retrieved records or update existing
func (repo *Repository) Upsert(sc domain.SDNRecord, updID string) error {
	_, err := repo.db.Exec(context.Background(), upsertSQL,
		sc.UID,
		sc.FirstName,
		sc.LastName,
	)

	return err
}

func (repo *Repository) QueryByName(name string) ([]domain.SDNRecord, error) {
	rows, err := repo.db.Query(context.Background(), queryStrictSQL, strings.ToLower(name))
	if err != nil {
		return nil, err
	}
	recs := make([]domain.SDNRecord, 0)
	for rows.Next() {
		rec := domain.SDNRecord{}
		err := rows.Scan(
			&rec.UID,
			&rec.FirstName,
			&rec.LastName,
			&rec.CreatedAt,
			&rec.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	return recs, nil
}

func (repo *Repository) QueryByPattern(query string) ([]domain.SDNRecord, error) {
	// create %pattern% wildcards array
	wildcards := make([]string, 0)

	for _, v := range strings.Split(query, " ") {
		wildcards = append(wildcards, "%"+v+"%")
	}

	rows, err := repo.db.Query(context.Background(), queryWeakSQL, wildcards)
	if err != nil {
		return nil, err
	}
	recs := make([]domain.SDNRecord, 0)
	for rows.Next() {
		rec := domain.SDNRecord{}
		err := rows.Scan(
			&rec.UID,
			&rec.FirstName,
			&rec.LastName,
			&rec.CreatedAt,
			&rec.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	return recs, nil
}
