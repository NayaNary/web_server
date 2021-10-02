package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"test.task/src/web"
)

func main() {

	webRouts := web.NewWebRouts()

	srv := &http.Server{
		Handler: webRouts.RoutConnect(),
		Addr:    "127.0.0.1:5000",
		WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
	}
	fmt.Println("server start: ",srv.Addr)
	log.Fatal(srv.ListenAndServe())

}