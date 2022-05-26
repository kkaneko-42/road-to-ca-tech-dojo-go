package handler

import (
	"log"
	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func HandleCollectionListGet(db *sql.DB) http.HandlerFunc {
	return func (writer http.ResponseWriter, req *http.Request) {
		items, err := getCollections(db, req)
		if err != nil {
			putError(writer, err)
			return;
		}

		jsondata, err := json.Marshal(&items)
		if err != nil {
			putError(writer, err)
		}
		writer.Write(jsondata)
	}
}

func getCollections(db *sql.DB, req *http.Request) ([]*CollectionGetResponce, error) {
	ctx := req.Context()
	user_id := ctx.Value("user_id")

	
}

type CollectionGetResponce struct {
	CollectionId string `json:"collectionID"`
	Name string `json:"name"`
	Rarity int `json:"rarity"`
	HasItem bool `json:"hasItem"`
}
