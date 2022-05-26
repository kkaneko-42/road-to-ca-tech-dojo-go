package handler

import (
	"log"
	"net/http"
)

func putError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
}
