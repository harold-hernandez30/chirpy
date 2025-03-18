package main

import (
	"log"
	"net/http"
)

func handleError(res http.ResponseWriter, err error, statusCode int, message string) error {
	if err != nil {
		log.Println(message)
		res.WriteHeader(statusCode)
		return err
	}

	return nil
}