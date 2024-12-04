package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type Storage interface {
	CreateToken(string, string, string, time.Time) error
	RefreshToken(string, string,  string, time.Time) error
	ValidateToken(string, string) (bool, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=root sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) ValidateToken(user_id string, refresh_token string) (bool, error) {
	user_id_uuid, err := uuid.FromString(user_id)
	if err != nil {
		return false, fmt.Errorf("error parsing user_id UUID: %w", err)
	}

	query := `select refresh_token from tokens
			where user_id = $1`
	rows, err := s.db.Query(query, user_id_uuid)
	if err != nil {
		return false, fmt.Errorf("failed to update token: %w", err)
	}

	var token string

	for rows.Next() {
		err = rows.Scan(token)
	}
	if err != nil {
		return false, fmt.Errorf("failed to validate token: %w", err)
	}
	return token == refresh_token, nil
}

func (s *PostgresStore) RefreshToken(user_id string, refresh_token string, ip_address string, expires_in time.Time) error {
	user_id_uuid, err := uuid.FromString(user_id)
	if err != nil {
		return fmt.Errorf("error parsing user_id UUID: %w", err)
	}

	query := `update tokens
		set refresh_token = $1
		where user_id = $2`

	_, err = s.db.Exec(query, refresh_token, user_id_uuid)
	if err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}
	return nil
}

func (s *PostgresStore) CreateToken(
		user_id string,
		refresh_token string,
		ip_address string,
		expires_in time.Time,
	) error {
	user_id_uuid, err := uuid.FromString(user_id)
	if err != nil {
		return fmt.Errorf("error parsing user_id UUID: %w", err)
	}

	query := `insert into tokens
	    (user_id, refresh_token, ip_address, email, expires_in)
		values ($1, $2, $3, $4)`

	_, err = s.db.Exec(query, user_id_uuid, refresh_token, ip_address, expires_in)
	if err != nil {
		return fmt.Errorf("failed to insert new token: %w", err)
	}
	return nil
}

func (s *PostgresStore) Init() error {
	return s.createTokenTable()
}

func (s *PostgresStore) createTokenTable() error {
	query := `create table if not exists tokens (
		user_id uuid NOT NULL, 
		refresh_token text NOT NULL,
		ip_address varchar(15) NOT NULL,
		expires_in timestamp NOT NULL,
        created_at timestamp without time zone NOT NULL default current_timestamp
	)`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}
