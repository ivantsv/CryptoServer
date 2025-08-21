package main

import (
	"log"
	"net/http"
	"crypto_server/auth"
	"crypto_server/db"
	"crypto_server/cryptoCRUD"

	"github.com/go-chi/chi/v5"
)

func main() {
	userDB := db.NewUserDB()
	cryptoDB := db.NewCryptoDB()

	authService := auth.NewAuthService(userDB)
	crudService := cryptocrud.NewCRUDService(cryptoDB)

	router := chi.NewRouter()
	router.Post("/auth/register", auth.RegisterHandler(authService))
	router.Post("/auth/login", auth.LoginHandler(authService))
	router.Get("/crypto", cryptocrud.GETHandlerCrypto(crudService))
	router.Post("/crypto", cryptocrud.POSTHandlerCrypto(crudService))

	log.Println("Starting server on http://localhost:8080")
	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatalln(err)
	}
}