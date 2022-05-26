package handler

import (
	"io"
	"encoding/json"
	"log"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// HandleUserUpdate ユーザー情報のupdate
func HandleUserUpdate(db *sql.DB) http.HandlerFunc {
	return func (writer http.ResponseWriter, req *http.Request) {
		err := updateUser(db, req)
		if err != nil {
			putError(writer, err)
			return
		}

		log.Print("User Update Successed")
	}
}

func updateUser(db *sql.DB, req *http.Request) error {
	ctx := req.Context()
	name_after, err := getRequestBody(req)
	if err != nil {
		return err
	}
	user_id := ctx.Value("user_id")

	_, err = db.Exec(
		"UPDATE users_infos " +
		"SET user_name = ? " +
		"WHERE user_id = ?", name_after, user_id)
	if err != nil {
		return err
	}
	return nil
}

func getRequestBody(req *http.Request) (string, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	var json_body userUpdateRequest
	err = json.Unmarshal(body, &json_body)
	if err != nil {
		return "", err
	}

	return json_body.Name, nil
}

type userUpdateRequest struct {
	Name string
}
