package storage

import (
	"database/sql"
	"time"

	"github.com/fatykhovar/jwtAuth/internal/domain"
	token_postgres "github.com/fatykhovar/jwtAuth/internal/storage/postgres/token"
	user_postgres "github.com/fatykhovar/jwtAuth/internal/storage/postgres/user"
	_ "github.com/lib/pq"
)

type User interface {
	CreateUser(string, string, string, string, time.Time) error
	GetUser(string) (domain.User, error)
}

type Token interface {
	RefreshToken(string, string, string, time.Time) error
}

type Storage struct {
	User
	Token
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		User:  user_postgres.NewUserPostgres(db),
		Token: token_postgres.NewTokenPostgres(db),
	}
}
