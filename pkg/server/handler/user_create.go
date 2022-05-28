package handler

import (
	"io"
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	id_len int = 8
	max_name_len int = 16
	token_len int = 16
	letters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// HandleUserCreate 新しいユーザーの追加
func HandleUserCreate(db *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		token, err, status := postUserData(db, request);

		if err != nil {
			putError(writer, err, status)
			return;
		}

		data, err := json.Marshal(&userCreateResponse{
			Token: token,
		})
		if err != nil {
			putError(writer, err, http.StatusInternalServerError)
			return
		}
		log.Print("User Creation Successed")
		writer.Write(data)
	}
}

func postUserData(db *sql.DB, req *http.Request) (string, error, int) {
	user, err := createUserData(req)
	if err != nil {
		return "", err, http.StatusInternalServerError
	}
	if len(string(user.name)) > max_name_len {
		log.Print("Name is too long")
		return "", fmt.Errorf("Name is too long"), http.StatusBadRequest
	}

	_, err = db.Exec("INSERT INTO users_tokens VALUES (?, ?)", user.id, user.token)
	if err != nil {
		return "", err, http.StatusInternalServerError
	}
	_, err = db.Exec("INSERT INTO users_infos VALUES (?, ?, 0)", user.id, user.name)
	if err != nil {
		return "", err, http.StatusInternalServerError
	}

	return user.token, nil, http.StatusOK
}

func createUserData(req *http.Request) (*userData, error) {
	id, err := generateRandomString(id_len)
	if err != nil {
		return nil, err
	}
	name, err := getUserName(req)
	if err != nil {
		return nil, err
	}
	token, err := generateRandomString(token_len)
	if err != nil {
		return nil, err
	}

	return &userData{
		id: id,
		name: name,
		token: token,
	}, nil
}

func getUserName(req *http.Request) (string, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	var json_body userCreateRequest
	err = json.Unmarshal(body, &json_body)
	if err != nil {
		return "", err
	}

	return json_body.Name, nil
}

type userCreateResponse struct {
	Token string `json:"token"`
}

type userCreateRequest struct {
	Name string
}

type userData struct {
	id string
	name string
	token string
}
