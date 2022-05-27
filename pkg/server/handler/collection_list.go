package handler

import (
	"net/http"
	"encoding/json"
	"database/sql"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis"
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

	res, err := getAllItems(cli)
	if err != nil {
		return nil, err
	}

	err = setItemHaving(res, user_id, db)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func getAllItems(cli *redis.Client) (*CollectionListGetResponce, error) {
	const nb_params int = 4
	var (
		items []Item
		buf []string
		rarity int
	)

	// get all keys in redis db
	keys, err := cli.Keys("*").Result()
	if err != nil {
		return nil, err
	}

	// get values and set into the struct
	for _, key := range keys {
		buf, err = cli.LRange(key, 0, -1).Result()
		if err != nil {
			return nil, err
		}

		rarity, _ = strconv.Atoi(buf[1])
		items = append(items, Item{
			CollectionId: key,
			Name: buf[0],
			Rarity: rarity,
			HasItem: false,
		})
	}

	return &CollectionListGetResponce{Collections: items}, nil
}

func setItemHaving(res *CollectionListGetResponce, user_id string, db *sql.DB) error {
	inventories, err := getUserInventories(db, user_id)
	if err != nil {
		return err
	}

	for i, item := range res.Collections {
		if strcontains(inventories, item.CollectionId) {
			res.Collections[i].HasItem = true
		}
	}
	return nil
}

func getUserInventories(db *sql.DB, user_id string) (*[]string, error) {
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
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &inventories, nil
}

type Item struct {
	CollectionId string `json:"collectionID"`
	Name string `json:"name"`
	Rarity int `json:"rarity"`
	HasItem bool `json:"hasItem"`
}

type CollectionListGetResponce struct {
	Collections []Item `json:"collections"`
}
