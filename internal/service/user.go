package service

import "github.com/fatykhovar/jwtAuth/internal/storage"

type UserService struct {
	store storage.User
}

func NewUserService(store storage.User) *UserService {
	return &UserService{store: store}
}
