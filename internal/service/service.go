package service

import "github.com/fatykhovar/jwtAuth/internal/storage"

type User interface {
	CreateUser(user user.User) (int, error)
}

type Token interface{

}


type Service struct {
	User
	Token
}

func NewService(store *storage.Storage) *Service {
	return &Service{
		User:	NewUserService(store.User),
		Token:	NewTokenService(store.Token),
	}
}
