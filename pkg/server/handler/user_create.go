package handler

import (
	"io"
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"crypto/rand"
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
		token, err := postUserData(db, request);

		if err != nil {
			log.Print("User Creation Failed: ", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return;
		}

		data, err := json.Marshal(&userCreateResponse{
			Token: token,
		})
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Print("User Creation Successed")
		writer.Write(data)
	}
}

func postUserData(db *sql.DB, req *http.Request) (string, error) {
	user, err := createUserData(req)
	if err != nil {
		return "", err
	}
	if len(string(user.name)) > max_name_len {
		log.Print("Name is too long")
		return "", fmt.Errorf("Name is too long")
	}

	_, err = db.Exec("INSERT INTO users_tokens VALUES (?, ?)", user.id, user.token)
	if err != nil {
		return "", err
	}
	_, err = db.Exec("INSERT INTO users_infos VALUES (?, ?, 0)", user.id, user.name)
	if err != nil {
		return "", err
	}

	return user.token, nil
}

func createUserData(req *http.Request) (*userData, error) {
	id, err_id := generateRandomString(id_len)
	if err_id != nil {
		return nil, err_id
	}
	name, err_name := getUserName(req)
	if err_name != nil {
		return nil, err_name
	}
	token, err_token := generateRandomString(token_len)
	if err_token != nil {
		return nil, err_token
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

func generateRandomString(length int) (string, error) {
	buf := make([]byte, length)
	var res string

	_, err := rand.Read(buf);
	if err != nil {
		return "", err
	}

	for _, v := range buf {
		res += string(letters[int(v) % len(letters)])
	}
	return res, nil
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
