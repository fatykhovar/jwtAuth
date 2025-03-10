package token_postgres

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type TokenPostgres struct {
	db *sql.DB
}

func NewTokenPostgres(db *sql.DB) *TokenPostgres {
	return &TokenPostgres{db: db}
}

func (s *TokenPostgres) RefreshToken(user_id string, refresh_token string, ip_address string, expires_in time.Time) error {
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