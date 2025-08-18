package auth

import (
	"crypto_server/db"
	"encoding/json"
	"log"
	"net/http"
)

type RegisterJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
    Token string `json:"token"`
}

type AuthService struct {
	UsersDB *db.UserDB
}

func NewAuthService(udb *db.UserDB) *AuthService {
	return &AuthService{UsersDB: udb}
}

func (authService *AuthService) Insert(login string, password string) error {
	return authService.UsersDB.Insert(login, password)
}

func (authService *AuthService) Exist(login string) bool {
	_, err := authService.UsersDB.Get(login)
	
	return err == nil
}

func (authService *AuthService) Delete(login string) error {
	return authService.UsersDB.Delete(login)
}

func RegisterHandler(authService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
			return
		}

		var registerJSON RegisterJSON

		err := json.NewDecoder(r.Body).Decode(&registerJSON)
		if err != nil {
			http.Error(w, `Bad Request`, http.StatusBadRequest)
			return
		}

		login := registerJSON.Username
		password := registerJSON.Password

		if authService.Exist(login) {
			http.Error(w, `User already exists`, http.StatusConflict)
			return
		}

		hashed_password := HashPassword(password)
		err = authService.Insert(login, hashed_password)
		if err != nil {
			log.Println("DB update error:", err)
			http.Error(w, `Server error`, http.StatusInternalServerError)
			return
		}

		tokenString, err := GenerateToken(login)
		if err != nil {
			log.Println("Token creation error:", err)
			http.Error(w, `Server error`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(TokenResponse{Token: tokenString})
	}	
}

func LoginHandler(authService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
		}

		
	}	
}