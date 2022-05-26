package middleware

import (
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
		user_id := ConfirmToken(request)
		if user_id == "" {
			log.Print("Auth Failed\n")
			writer.WriteHeader(http.StatusInternalServerError)
			return;
		} else {
			ctx = context.WithValue(ctx, "user_id", user_id)
		}

		nextFunc(writer, request.WithContext(ctx))
	}
}

func ConfirmToken(req *http.Request) string {

	/* db接続 */
	db, err := sql.Open("mysql", "root:ca-tech-dojo@(127.0.0.1:3306)/road_to_ca")
	if err != nil {
		log.Fatal("DB connection failed: ", err)
		return ""
	}
	defer db.Close()

	/* tokenの照合 */
	token := req.URL.Query().Get("x-token")
	var user_id string

	err = db.QueryRow("SELECT user_id FROM users_tokens WHERE token = ?;", token).Scan(&user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		} else {
			log.Fatal("DB Query Error: ", err)
		}
	}
	return user_id
}
