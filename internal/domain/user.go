package domain

import "time"

type User struct {
	UserID       string
	RefreshToken string
	Email        string
	IpAddress    string
	ExpiresIn    time.Time
}