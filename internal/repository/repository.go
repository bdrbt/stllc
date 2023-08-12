package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bdrbt/stllc/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
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

func New(dbURL string) (*Repository, error) {
	repo := &Repository{}

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_URL:", err)
	}

	repo.db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}

	return repo, nil
}

// Upsert - insert retrieved records or update existing.
func (repo *Repository) Upsert(ctx context.Context, sc domain.SDNRecord) error {
	_, err := repo.db.Exec(ctx, upsertSQL,
		sc.UID,
		sc.FirstName,
		sc.LastName,
	)

	return fmt.Errorf("error upserting record:%w", err)
}

func (repo *Repository) QueryByName(ctx context.Context, name string) ([]domain.SDNRecord, error) {
	rows, err := repo.db.Query(ctx, queryStrictSQL, strings.ToLower(name))
	if err != nil {
		return nil, fmt.Errorf("error querying records:%w", err)
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
			return nil, fmt.Errorf("error mapping values:%w", err)
		}

		recs = append(recs, rec)
	}

	return recs, nil
}

func (repo *Repository) QueryByPattern(ctx context.Context, query string) ([]domain.SDNRecord, error) {
	wildcards := make([]string, 0)

	// fill %pattern% wildcards array
	for _, v := range strings.Split(query, " ") {
		wildcards = append(wildcards, "%"+v+"%")
	}

	rows, err := repo.db.Query(ctx, queryWeakSQL, wildcards)
	if err != nil {
		return nil, fmt.Errorf("error querying records:%w", err)
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
			return nil, fmt.Errorf("error mapping values:%w", err)
		}

		recs = append(recs, rec)
	}

	return recs, nil
}
