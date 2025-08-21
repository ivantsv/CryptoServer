package auth

import (
	"crypto_server/db"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var ErrInvalidUserData = errors.New("incorrect user data")

type LoginPasswordJSON struct {
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

func (authService *AuthService) UserValidation(login string, password string) error {
	realPassword, err := authService.UsersDB.Get(login)
	if err != nil {
		return err
	}

	if DifferentPasswords(password, realPassword) {
		return ErrInvalidUserData
	}

	return nil
}

func (authService *AuthService) Delete(login string) error {
	return authService.UsersDB.Delete(login)
}

func RegisterHandler(authService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var LoginPasswordJSON LoginPasswordJSON

		err := json.NewDecoder(r.Body).Decode(&LoginPasswordJSON)
		if err != nil {
			http.Error(w, `Bad Request`, http.StatusBadRequest)
			return
		}

		login := LoginPasswordJSON.Username
		password := LoginPasswordJSON.Password

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
		var loginJSON LoginPasswordJSON

		err := json.NewDecoder(r.Body).Decode(&loginJSON)
		if err != nil {
			log.Println("Error during JSON parsing: ", err)
			http.Error(w, `Bad Request`, http.StatusBadRequest)
			return
		}

		login := loginJSON.Username
		password := loginJSON.Password

		err = authService.UserValidation(login, HashPassword(password))
		if err != nil {
			log.Println("Error during user validation: ", err)
			http.Error(w, `Incorrect login or password`, http.StatusUnauthorized)
			return 
		}

		tokenString, err := GenerateToken(login)
		if err != nil {
			log.Println("Token creation error:", err)
			http.Error(w, `Server error`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(TokenResponse{Token: tokenString})
	}	
}