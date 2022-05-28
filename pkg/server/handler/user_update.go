package handler

import (
	"io"
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// HandleUserUpdate ユーザー情報のupdate
func HandleUserUpdate(db *sql.DB) http.HandlerFunc {
	return func (writer http.ResponseWriter, req *http.Request) {
		err, status := updateUser(db, req)
		if err != nil {
			putError(writer, err, status)
			return
		}

		log.Print("User Update Successed")
	}
}

func updateUser(db *sql.DB, req *http.Request) (error, int) {
	name_after, err := getRequestBody(req)
	if err != nil {
		return err, http.StatusInternalServerError
	} else if len(name_after) > max_name_len {
		return fmt.Errorf("Name is too long"), http.StatusBadRequest
	}
	user_id := getUserIdFromContext(req)

	_, err = db.Exec(
		"UPDATE users_infos " +
		"SET user_name = ? " +
		"WHERE user_id = ?", name_after, user_id)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
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
