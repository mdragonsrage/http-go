package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mdragonsrage/http-go/house"
)

func main() {

	s := house.NewInMemoryStorage()
	router := http.NewServeMux()

	router.HandleFunc("/healthcheck", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Contet-Type", "application-json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{"message": "service is running"}`))
	})

	router.HandleFunc("/lightbulbs", house.GetLightbulb(s))
	router.HandleFunc("/lightbulbs/switch", house.SwitchLightbulb(s))
	router.HandleFunc("/lightbulbs/create", house.CreateLightbulb(s))
	router.HandleFunc("/lightbulbs/delete", house.DeleteLightbulb(s))

	srv := http.Server{
		Addr:         ":8080",
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
		Handler:      router,
	}

	fmt.Println("http server listing on localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
