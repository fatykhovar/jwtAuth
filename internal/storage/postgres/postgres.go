package postgres

import (
	"database/sql"
	"fmt"

	"github.com/fatykhovar/jwtAuth/internal/config"
)

func NewPostgresDB(cfg config.Config) (*sql.DB, error) {
	const op = "postgres.NewPostgresDB"

	connStr := fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Storage.Host,
		cfg.Storage.User,
		cfg.Storage.DBname,
		cfg.Storage.Password,
		cfg.Storage.SSLMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		create table if not exists user (
			user_id uuid PRIMARY KEY, 
			ip_address varchar(15) NOT NULL,
			email varchar(255) NOT NULL);

		create table if not exists token (
			id integer PRIMARY KEY,
			user_id uuid PRIMARY KEY, 
			refresh_token text NOT NULL,
			expires_in timestamp NOT NULL);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}