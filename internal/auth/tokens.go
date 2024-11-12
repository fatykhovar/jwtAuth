package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fatykhovar/jwtAuth/internal/storage"
	jwt "github.com/golang-jwt/jwt/v4"

	"github.com/gorilla/mux"
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

// Настройки JWT
var (
	secretKey       = []byte("secret-key")
	accessTokenTTL  = time.Hour * 1
	refreshTokenTTL = time.Hour * 24
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
	jwtToken   Token
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/token/{userID}", makeHTTPHandleFunc(s.handleNewToken)).Methods("GET")
	router.HandleFunc("/token/{userID}", makeHTTPHandleFunc(s.handleRefreshToken)).Methods("POST")

	http.ListenAndServe(s.listenAddr, router)
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Генерация JWT токена
func generateToken(guid string, ttl time.Duration, ip string) (string, error) {
	claims := jwt.MapClaims{
		"guid":      guid,
		"expiresAt": time.Now().Add(ttl),
		"ip":        ip,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func (s *APIServer) handleNewToken(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID := vars["userID"]

	if userID == "" {
		return fmt.Errorf("UserID is required", http.StatusBadRequest)
	}

	ip := ip_address[0]

	// Генерация Access и Refresh токенов
	accessToken, err := generateToken(userID, accessTokenTTL, ip)
	if err != nil {
		return fmt.Errorf(err.Error(), http.StatusInternalServerError)
	}

	refreshToken, err := generateToken(userID, refreshTokenTTL, ip)
	if err != nil {
		return fmt.Errorf(err.Error(), http.StatusInternalServerError)
	}

	tokens := Token{AccessToken: accessToken, RefreshToken: refreshToken}
	s.store.CreateToken(userID, refreshToken, ip, time.Now().Add(refreshTokenTTL))
	return WriteJSON(w, http.StatusOK, tokens)
}

func (s *APIServer) handleRefreshToken(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID := vars["userID"]

	if userID == "" {
		return fmt.Errorf("UserID is required", http.StatusBadRequest)
	}

	ip := ip_address[0]

	// Генерация Access и Refresh токенов
	accessToken, err := generateToken(userID, accessTokenTTL, ip)
	if err != nil {
		return fmt.Errorf(err.Error(), http.StatusInternalServerError)
	}

	refreshToken, err := generateToken(userID, refreshTokenTTL, ip)
	if err != nil {
		return fmt.Errorf(err.Error(), http.StatusInternalServerError)
	}

	tokens := Token{AccessToken: accessToken, RefreshToken: refreshToken}
	s.store.RefreshToken(userID, refreshToken, ip, time.Now().Add(refreshTokenTTL))
	return WriteJSON(w, http.StatusOK, tokens)
}
