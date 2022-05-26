package handler

import (
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
			return
		}

		jsondata, err := json.Marshal(&items)
		if err != nil {
			putError(writer, err)
			return
		}
		writer.Write(jsondata)
	}
}

func getCollections(db *sql.DB, req *http.Request) (*CollectionListGetResponce, error) {
	rows, err := db.Query(
		"SELECT item_id, item_name, rarity, user_id FROM items " + 
		"LEFT OUTER JOIN users_inventories " +
		"USING (item_id);")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := parseDBReturn(rows, req)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func parseDBReturn (rows *sql.Rows, req *http.Request) (*CollectionListGetResponce, error) {
	var res CollectionListGetResponce
	var row returnedRow
	user_id := getUserIdFromContext(req)

	for rows.Next() {
		err := rows.Scan(
			&row.item_id,
			&row.item_name,
			&row.rarity,
			&row.having_user_id)
		if err != nil {
			return nil, err
		}

		res.Collections = append(res.Collections, &Item{
			CollectionId: row.item_id,
			Name: row.item_name,
			Rarity: row.rarity,
			HasItem: checkItemHaving(&row.having_user_id, user_id),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &res, nil
}

func checkItemHaving (having_user_id *sql.NullString, user_id string) bool {

	if having_user_id.Valid && having_user_id.String == user_id {
		return true
	} else {
		return false
	}
}

type Item struct {
	CollectionId string `json:"collectionID"`
	Name string `json:"name"`
	Rarity int `json:"rarity"`
	HasItem bool `json:"hasItem"`
}

type returnedRow struct {
	item_id string
	item_name string
	rarity int
	having_user_id sql.NullString
}

type CollectionListGetResponce struct {
	Collections []*Item `json:"collections"`
}
