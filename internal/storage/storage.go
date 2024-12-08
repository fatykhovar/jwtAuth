package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type PostgresStore struct {
	db *sql.DB
}

type Storage interface {
	CreateUser(string, string, string, string, time.Time) error
	RefreshToken(string, string,  string, time.Time) error
	GetUser(string) (User, error)
}

type User struct{
	UserID string
	RefreshToken string
	Email string
	IpAddress string
	ExpiresIn time.Time
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

func (s *PostgresStore) GetUser(user_id string) (User, error) {
	user_id_uuid, err := uuid.FromString(user_id)
	if err != nil {
		return User{}, fmt.Errorf("error parsing user_id UUID: %w", err)
	}

	query := `select * from tokens
			where user_id = $1`
	rows, err := s.db.Query(query, user_id_uuid)
	if err != nil {
		return User{}, err
	}

	var user User

	for rows.Next() {
		err = rows.Scan(&user.UserID,
						&user.RefreshToken,
					    &user.IpAddress,
						&user.Email,
                        &user.ExpiresIn)
	}
	if err != nil {
		return User{}, err
	}
	return user, nil
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

func (s *PostgresStore) CreateUser(
		user_id string,
		refresh_token string,
		ip_address string,
		email string,
		expires_in time.Time,
	) error {
	user_id_uuid, err := uuid.FromString(user_id)
	if err != nil {
		return fmt.Errorf("error parsing user_id UUID: %w", err)
	}

	query := `insert into tokens
	    (user_id, refresh_token, ip_address, email, expires_in)
		values ($1, $2, $3, $4, $5)`

	_, err = s.db.Exec(query, user_id_uuid, refresh_token, ip_address, email, expires_in)
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
		email varchar(255) NOT NULL,
		expires_in timestamp NOT NULL
	)`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}
