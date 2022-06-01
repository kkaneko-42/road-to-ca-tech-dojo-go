package middleware

import (
	"fmt"
	"log"
	"context"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID string `db:"user_id"`
	Token string `db:"token"`
}

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		ctx := request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		// TODO: implement here
		user_id, err := ConfirmToken(request)
		if err != nil {
			log.Println("Auth Failed: ", err)
			writer.WriteHeader(http.StatusUnauthorized)
			return;
		} else {
			ctx = context.WithValue(ctx, "user_id", user_id)
		}
		log.Println("Auth successed")
		nextFunc(writer, request.WithContext(ctx))
	}
}

func ConfirmToken(req *http.Request) (string, error) {

	/* db接続 */
	db, err := sql.Open("mysql", "root:ca-tech-dojo@(mysql:3306)/road_to_ca")
	if err != nil {
		return "", err
	}
	defer db.Close()

	/* tokenの照合 */
	token := req.Header.Get("X-token")
	var user_id string

	/* tokenで検索し、検索結果なしなら認証失敗 */
	err = db.QueryRow("SELECT user_id FROM users_tokens WHERE token = ?;", token).Scan(&user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("User not found")
		} else {
			return "", err
		}
	}
	return user_id, nil
}
