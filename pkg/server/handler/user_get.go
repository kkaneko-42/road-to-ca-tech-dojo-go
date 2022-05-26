package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// HandleUserCreate 新しいユーザーの追加
func HandleUserGet(db *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		user_data, err := getUserData(db, request)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsondata, err := json.Marshal(&user_data)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Print("User Get Successed")
		writer.Write(jsondata)
	}
}

func getUserData(db *sql.DB, req *http.Request) (*userGetResponse, error) {
	var user_data userGetResponse
	ctx := req.Context()

	user_data.Id = ctx.Value("user_id").(string)
	err := db.QueryRow(
		"SELECT user_name, having_coins, max_score " +
		"FROM users_infos " +
		"INNER JOIN (" +
			"SELECT user_id, MAX(score) AS max_score FROM scores " +
			"GROUP BY (user_id)) " +
		"AS max_scores " +
		"USING (user_id) " + 
		"WHERE user_id = ?", user_data.Id).Scan(
			&(user_data.Name),
			&(user_data.HavingCoins),
			&(user_data.HighScore))
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
