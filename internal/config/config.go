package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Postgres struct {
		Host         string
		Port         int
		User         string
		Password     string
		Database     string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Addr string
}

func Load() (*Config, error) {
	var err error

	cfg := &Config{}
	cfg.Postgres.Host = os.Getenv("PG_HOST")
	cfg.Postgres.User = os.Getenv("PG_USER")
	cfg.Postgres.Password = os.Getenv("PG_PASS")
	cfg.Postgres.Database = os.Getenv("PG_DATABASE")

	if cfg.Postgres.MaxIdleTime = os.Getenv("PG_IDLE_TIME"); cfg.Postgres.MaxIdleTime == "" {
		cfg.Postgres.MaxIdleTime = "10m"
	}

	if cfg.Postgres.Port, err = strconv.Atoi(os.Getenv("PG_PORT")); err != nil {
		cfg.Postgres.Port = 5432
	}

	if cfg.Postgres.MaxOpenConns, err = strconv.Atoi(os.Getenv("PG_OPEN_CONNS")); err != nil {
		cfg.Postgres.MaxOpenConns = 4
	}

	if cfg.Postgres.MaxIdleConns, err = strconv.Atoi(os.Getenv("PG_IDLE_CONNS")); err != nil {
		cfg.Postgres.MaxIdleConns = 4
	}

	if cfg.Addr = os.Getenv("ADDR"); cfg.Addr == "" {
		cfg.Addr = ":8080"
	}

	return cfg, cfg.Validate()
}

//nolint:goerr113
func (cfg *Config) Validate() error {
	if cfg.Postgres.Host == "" {
		return errors.New("postgres host is not set")
	}

	if cfg.Postgres.User == "" || cfg.Postgres.Password == "" {
		return errors.New("postgrs username/password is not set")
	}

	if cfg.Postgres.Database == "" {
		return errors.New("database name is not set")
	}

	return nil
}

//nolint:nosprintfhostport
func (cfg *Config) PgURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)
}
