package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"github.com/fatykhovar/jwtAuth/internal/storage"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}

type Token struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	IPAddress    string `json:"ip_address"`
}

type ApiError struct {
	Error string `json:"error"`
}

// Настройки JWT
var (
	secretKey       = []byte("secret-key")
	accessTokenTTL  = time.Hour * 1
	refreshTokenTTL = time.Hour * 24
)
type TokenService struct {
	store storage.Token
}

func NewTokenService(store storage.Token) *TokenService {
	return &TokenService{store: store}
}

func generateAccessToken(guid string, ttl time.Duration, ip string) (string, error) {
	claims := jwt.MapClaims{
		"guid":      guid,
		"expiresAt": time.Now().Add(ttl),
		"ip":        ip,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func generateRefreshToken(ttl time.Duration, ip string) (string, error) {
	claims := jwt.MapClaims{
		"expiresAt": time.Now().Add(ttl),
		"ip":        ip,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func (s *TokenService) GetToken(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID := vars["userID"]

	if userID == "" {
		return fmt.Errorf("userID is required", http.StatusBadRequest)
	}

	ip := IP_address[0]

	user, err := s.store.GetUser(userID)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	// Генерация Access и Refresh токенов
	accessToken, err := generateAccessToken(userID, accessTokenTTL, ip)
	if err != nil {
		return fmt.Errorf(err.Error(), http.StatusInternalServerError)
	}
	tokens := Token{AccessToken: accessToken,
			RefreshToken: user.RefreshToken}
	return WriteJSON(w, http.StatusOK, tokens)
}

func (s *TokenService) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID := vars["userID"]

	if userID == "" {
		return fmt.Errorf("userID is required", http.StatusBadRequest)
	}

	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("error decoding request body: %w", err)
	}
	defer r.Body.Close()

	user, err := s.store.GetUser(userID)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	hashedRefreshToken := user.RefreshToken

	// валидация Refresh токена
	err = bcrypt.CompareHashAndPassword([]byte(hashedRefreshToken), []byte(req.RefreshToken))
	if err != nil {
		return WriteJSON(w, http.StatusForbidden, ApiError{Error: "invalid refresh token"})
	}

	// парсинг claims
	type RefreshTokenClaims struct {
		expiresAt time.Duration
		ip string
		jwt.RegisteredClaims
	} 
	token, err := jwt.ParseWithClaims(req.RefreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return WriteJSON(w, http.StatusForbidden, ApiError{Error: "failed to get claims"})
	}

	if err := sendEmail(Email); err != nil {
		fmt.Errorf(err.Error(), http.StatusInternalServerError)
	}
	// валидация ip
	if claims.ip != user.IpAddress {
		if err := sendEmail(Email); err != nil {
			fmt.Errorf(err.Error(), http.StatusInternalServerError)
		}
	}
	// Генерация Access и Refresh токенов
	accessToken, err := generateAccessToken(userID, accessTokenTTL, req.IPAddress)
	if err != nil {
		return fmt.Errorf(err.Error(), http.StatusInternalServerError)
	}

	refreshToken, err := generateRefreshToken(refreshTokenTTL, req.IPAddress)
	if err != nil {
		return fmt.Errorf(err.Error(), http.StatusInternalServerError)
	}

	tokens := Token{AccessToken: accessToken, RefreshToken: refreshToken}
	s.store.RefreshToken(userID, refreshToken, req.IPAddress, time.Now().Add(refreshTokenTTL))
	return WriteJSON(w, http.StatusOK, tokens)
}

func sendEmail(email string) error {
	from := "test"
	password := "test"

	to := []string{
	   email,
	}
 
	// smtp сервер конфигурация
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
 
	message := []byte("IP-адрес не действителен.")

	auth := smtp.PlainAuth("", from, password, smtpHost)
 
	err := smtp.SendMail(smtpHost + ":" + smtpPort, auth, from, to, message)
	if err != nil {
	   return fmt.Errorf("error sending email: %w", err)
	}
	return nil
}

func (s *APIServer) InsertTestData(ip string) error{
	userID := uuid.New()

	refreshToken, err := generateRefreshToken(refreshTokenTTL, ip)
	if err != nil {
		return fmt.Errorf(err.Error(), http.StatusInternalServerError)
	}

	encodedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf(err.Error(), http.StatusInternalServerError)
	}

	s.store.CreateUser(userID.String(), string(encodedRefreshToken), ip, Email, time.Now().Add(refreshTokenTTL))
	return nil
}