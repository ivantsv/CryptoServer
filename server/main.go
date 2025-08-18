package main

import (
	"log"
	"net/http"
	"crypto_server/auth"
	"crypto_server/db"
)

func main() {
	userDB := db.NewUserDB()
	authService := auth.NewAuthService(userDB)

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/auth/register", auth.RegisterHandler(authService))

	log.Println("Starting server on http://localhost:8080")
	err := http.ListenAndServe("localhost:8080", serverMux)
	if err != nil {
		log.Fatalln(err)
	}
}