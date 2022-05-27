package handler

import (
	"log"
	"net/http"
	"crypto/rand"
)

func putError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
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

func strcontains(arr *[]string, str string) bool {
	for _, v := range *arr {
		if v == str {
			return true
		}
	}
	return false
}
