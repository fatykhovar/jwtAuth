package user_postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fatykhovar/jwtAuth/internal/domain"
	uuid "github.com/satori/go.uuid"
)

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (s *UserPostgres) GetUser(user_id string) (domain.User, error) {
	const op = "storage.GetUser"

	user_id_uuid, err := uuid.FromString(user_id)
	if err != nil {
		return domain.User{}, fmt.Errorf("error parsing user_id UUID: %w", err)
	}

	query := `select * from tokens
			where user_id = $1`
	rows, err := s.db.Query(query, user_id_uuid)
	if err != nil {
		return domain.User{}, err
	}

	var user domain.User

	for rows.Next() {
		err = rows.Scan(&user.UserID,
			&user.RefreshToken,
			&user.IpAddress,
			&user.Email,
			&user.ExpiresIn)
	}
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *UserPostgres) CreateUser(
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