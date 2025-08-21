package cryptocrud

import (
	"crypto_server/db"
	"encoding/json"
	"net/http"
)

type CRUDService struct {
	cryptoDB *db.CryptoDB
}

func NewCRUDService(cdb *db.CryptoDB) *CRUDService {
	return &CRUDService{cryptoDB: cdb}
}

func (crudService *CRUDService) AddNewCrypto(symbol string, coinData db.CoinData) error {
	return crudService.cryptoDB.Insert(symbol, coinData)
}

type CryptoResponse struct {
	Symbol string `json:"symbol"`
	db.CoinData
}

type AllCryptosResponse struct {
	Cryptos []CryptoResponse `json:"cryptos"`
}

type OneCryptoResponse struct {
	Crypto CryptoResponse `json:"crypto"`
}

func GETHandlerCrypto(crudSerivce *CRUDService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cryptos AllCryptosResponse
		for k, v := range crudSerivce.cryptoDB.Storage {
			cryptos.Cryptos = append(cryptos.Cryptos, CryptoResponse{Symbol: k, CoinData: v})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cryptos)
	}
}

// @TODO Запрос идет чисто symbol. Парсинг информации о криптовалюте через API CoinGecko
func POSTHandlerCrypto(crudSerivce *CRUDService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cryptoResponse CryptoResponse

		err := json.NewDecoder(r.Body).Decode(&cryptoResponse)
		if err != nil {
			http.Error(w, `Bad Request`, http.StatusBadRequest)
			return
		}

		err = crudSerivce.AddNewCrypto(cryptoResponse.Symbol, cryptoResponse.CoinData)
		if err != nil {
			http.Error(w, `Server error`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type",  "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(OneCryptoResponse{Crypto: cryptoResponse})
	}
}