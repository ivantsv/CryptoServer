package db

import (
	"errors"
	"sync"
	"time"
)

var ErrUnknownCoin = errors.New("Unknown coin name")

type CoinData struct {
	Name string `json:"name"`
	CurrentPrice float64 `json:"current_price"`
	LastUpdate time.Time `json:"last_updated"`
}

type CryptoDB struct {
	storage map[string]CoinData
	mutex sync.Mutex 
}

func NewCryptoDB() *CryptoDB {
	return &CryptoDB{
		storage: make(map[string]CoinData),
	}
}

func (cdb *CryptoDB) Insert(name string, data CoinData) error {
	cdb.mutex.Lock()
	defer cdb.mutex.Unlock()
	cdb.storage[name] = data
	return nil
}

func (cdb *CryptoDB) Get(name string) (CoinData, error) {
	cdb.mutex.Lock()
	defer cdb.mutex.Unlock()
	coinData, ok := cdb.storage[name]
	if !ok {
		return CoinData{}, ErrUnknownCoin
	}

	return coinData, nil
}

func (cdb *CryptoDB) Delete(name string) error {
	cdb.mutex.Lock()
	defer cdb.mutex.Unlock()
	_, ok := cdb.storage[name]
	if !ok {
		return ErrUnknownCoin
	}

	delete(cdb.storage, name)
	return nil
}