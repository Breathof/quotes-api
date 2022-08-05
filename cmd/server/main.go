package main

import (
	"awesomeProject/pkg/api"
	"awesomeProject/pkg/db"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	router := api.NewRouter()
	port := os.Getenv("PORT")
	log.Print("Init success.")

	_, err := db.NewDB()
	if err != nil {
		panic(err)
	}
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Printf("error on router init: %v\n", err)
	}
}
