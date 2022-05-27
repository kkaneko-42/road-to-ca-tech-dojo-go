package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/go-redis/redis"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"42tokyo-road-to-dojo-go/pkg/server/cache"
)

// HandleUserGet ユーザー情報の取得
func HandleUserGet(db *sql.DB, cli *redis.Client) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		user_data, err := getUserData(db, cli, request)
		if err != nil {
			putError(writer, err)
			return
		}

		jsondata, err := json.Marshal(&user_data)
		if err != nil {
			putError(writer, err)
			return
		}
		log.Print("User Get Successed")
		writer.Write(jsondata)
	}
}

func getUserData(db *sql.DB, cli *redis.Client, req *http.Request) (*userGetResponse, error) {
	var user_data userGetResponse

	user_data.Id = getUserIdFromContext(req)
	log.Println(user_data.Id)
	err := db.QueryRow(
		"SELECT user_name, having_coins FROM users_infos;").Scan(
			&(user_data.Name),
			&(user_data.HavingCoins))
	if err != nil {
		return nil, err
	}

	user_data.HighScore, err = cache.GetUserHighScore(user_data.Id, cli)
	if err != nil {
		return nil, err
	}

	return &user_data, nil
}

type userGetResponse struct {
	Id string `json:"id"`
	Name string `json:"name"`
	HighScore int `json:"highScore"`
	HavingCoins int `json:"coin"`
}
