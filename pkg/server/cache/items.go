package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	Item_key string = "items"
)

func CacheItems(db *sql.DB, cli *redis.Client) error {

	rows, err := db.Query("SELECT * FROM items;")
	if err != nil {
		return err
	}

	err = pushRows(rows, cli)
	if err != nil {
		return err
	}

	return nil
}

func pushRows(rows *sql.Rows, cli *redis.Client) error {
	var (
		err error
		item ItemData
		item_json []byte
	)

	for rows.Next() {
		err = rows.Scan(
			&item.Id,
			&item.Name,
			&item.Rarity,
			&item.GachaWeight,
		)
		if err != nil {
			return err
		}

		item_json, err = json.Marshal(&item)
		if err != nil {
			return err
		}

		err = cli.RPush(Item_key, item_json).Err()
		if err != nil {
			return err
		}
	}

	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func GetItems(cli *redis.Client) (*[]ItemData, error) {
	var (
		items []ItemData
		item ItemData
	)

	items_jsons, err := cli.LRange(Item_key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	for i, _ := range items_jsons {
		err = json.Unmarshal([]byte(items_jsons[i]), &item)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}
	return &items, nil
}

type ItemData struct {
	Id string
	Name string
	Rarity int
	GachaWeight int
}
