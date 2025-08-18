package db

import (
	"errors"
	"sync"
)

var (
	ErrLoginUsed   = errors.New("login already used")
	ErrUnknownUser = errors.New("unknown login")
)

type UserDB struct {
	loginPassword map[string]string
	mutex sync.Mutex
}

func NewUserDB() *UserDB {
    return &UserDB{
        loginPassword: make(map[string]string),
    }
}

func (udb *UserDB) Insert(login string, password string) error {
	udb.mutex.Lock()
	defer udb.mutex.Unlock()
	_, ok := udb.loginPassword[login]
	if ok {
		return ErrLoginUsed
	}

	udb.loginPassword[login] = password
	return nil
}

func (udb *UserDB) Get(login string) (string, error) {
	udb.mutex.Lock()
	defer udb.mutex.Unlock()
	password, ok := udb.loginPassword[login]
	if !ok {
		return "", ErrUnknownUser
	}

	return password, nil
}

func (udb *UserDB) Delete(login string) error {
	udb.mutex.Lock()
	defer udb.mutex.Unlock()
	_, ok := udb.loginPassword[login]
	if !ok {
		return ErrUnknownUser
	}

	delete(udb.loginPassword, login)
	return nil
}