package handler

import (
	"log"
	"net/http"
	"crypto/rand"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func putError(w http.ResponseWriter, err error, http_status int) {
	log.Println(err)
	w.WriteHeader(http_status)
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

func getUserIdFromContext(req *http.Request) string {
	return (req.Context().Value("user_id").(string))
}

func strcontains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func GetUserInventories(db *sql.DB, user_id string) ([]string, error) {
	var (
		inventories []string
		buf string
	)

	rows, err := db.Query(
		"SELECT item_id FROM users_inventories " + 
		"WHERE user_id = ?;", user_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&buf)
		if err != nil {
			return nil, err
		}
		inventories = append(inventories, buf)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inventories, nil
}
