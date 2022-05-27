package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// HandleUserGet ユーザー情報の取得
func HandleUserGet(db *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		user_data, err := getUserData(db, request)
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

func getUserData(db *sql.DB, req *http.Request) (*userGetResponse, error) {
	var (
		user_data userGetResponse
		highscore_buf sql.NullInt64
	)

	user_data.Id = getUserIdFromContext(req)
	log.Println(user_data.Id)
	err := db.QueryRow(
		"SELECT user_name, having_coins, max_score " +
		"FROM users_infos " +
		"LEFT OUTER JOIN (" +
			"SELECT user_id, MAX(score) AS max_score FROM scores " +
			"GROUP BY (user_id)) " +
		"AS max_scores " +
		"USING (user_id) " + 
		"WHERE user_id = ?;", user_data.Id).Scan(
			&(user_data.Name),
			&(user_data.HavingCoins),
			&(highscore_buf))
	if err != nil {
		return nil, err
	}

	user_data.HighScore = int(highscore_buf.Int64)
	return &user_data, nil
}

type userGetResponse struct {
	Id string `json:"id"`
	Name string `json:"name"`
	HighScore int `json:"highScore"`
	HavingCoins int `json:"coin"`
}
