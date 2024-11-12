package storage

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
)

type Storage interface {
	CreateToken(string, string, string, time.Time) (error)
	RefreshToken(string, string, string, time.Time) (error)
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

func (s *PostgresStore) RefreshToken(user_id string, refresh_token string, ip_address string, expires_in time.Time) error {
	user_id_uuid, err := uuid.FromString(user_id)
    if err != nil {
        log.Fatalf("Error parsing user_id UUID: %v", err)
	}

	query := `update tokens
		set refresh_token = $1
		where user_id = $2`

	_, err = s.db.Exec(query, refresh_token, user_id_uuid)
	if err != nil {
		log.Fatalf("Failed to update token: %v", err)
	}
	return err
}

func (s *PostgresStore) CreateToken(user_id string, refresh_token string, ip_address string, expires_in time.Time) error {
	user_id_uuid, err := uuid.FromString(user_id)
    if err != nil {
        log.Fatalf("Error parsing user_id UUID: %v", err)
	}

	query := `insert into tokens
	    (user_id, refresh_token, ip_address, expires_in)
		values ($1, $2, $3, $4)`

	_, err = s.db.Exec(query, user_id_uuid, refresh_token, ip_address, expires_in)
	if err != nil {
		log.Fatalf("Failed to insert new token: %v", err)
	}
	return err
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
		log.Fatalf("Failed to create table: %v", err)
	}
	return err
}
