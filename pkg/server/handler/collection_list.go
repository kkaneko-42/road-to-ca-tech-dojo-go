package handler

import (
	"log"
	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis"
	"42tokyo-road-to-dojo-go/pkg/server/cache"
)

func HandleCollectionListGet(db *sql.DB, cli *redis.Client) http.HandlerFunc {
	return func (writer http.ResponseWriter, req *http.Request) {
		res, err := createResponce(db, cli, req)
		if err != nil {
			putError(writer, err)
			return
		}

		jsondata, err := json.Marshal(&res)
		if err != nil {
			putError(writer, err)
			return
		}
		writer.Write(jsondata)
	}
}

func createResponce(db *sql.DB, cli *redis.Client, req *http.Request) (*CollectionListGetResponce, error) {
	user_id := getUserIdFromContext(req)

	items, err := cache.GetItems(cli)
	if err != nil {
		return nil, err
	}

	responce, err := setItemHaving(items, user_id, db)
	if err != nil {
		return nil, err
	}
	return responce, nil
}

func setItemHaving(items *[]cache.ItemData, user_id string, db *sql.DB) (*CollectionListGetResponce, error) {
	var (
		res CollectionListGetResponce
		item cache.ItemData
		has_item bool
	)

	inventories, err := getUserInventories(db, user_id)
	if err != nil {
		return nil, err
	}

	log.Println(*items)
	for i, _ := range *items {
		item = (*items)[i]
		if strcontains(inventories, item.Id) {
			has_item = true
		} else {
			has_item = false
		}

		res.Collections = append(res.Collections, &Item{
			CollectionId: item.Id,
			Name: item.Name,
			Rarity: item.Rarity,
			HasItem: has_item,
		})
	}

	return &res, nil
}

func getUserInventories(db *sql.DB, user_id string) ([]string, error) {
	var (
		inventories []string
		buf string
	)

	rows, err := db.Query(
		"SELECT item_id FROM users_inventories " + 
		"WHERE user_id = ?;", user_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&buf)
		if err != nil {
			return nil, err
		}
		inventories = append(inventories, buf)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inventories, nil
}

type Item struct {
	CollectionId string `json:"collectionID"`
	Name string `json:"name"`
	Rarity int `json:"rarity"`
	HasItem bool `json:"hasItem"`
}

type CollectionListGetResponce struct {
	Collections []*Item `json:"collections"`
}
