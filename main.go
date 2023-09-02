package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	lightbulbs = make(map[string]bool)
)

func main() {
	lightbulbs["livingroom"] = false
	lightbulbs["kitchen"] = false

	http.HandleFunc("/healthcheck", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Contet-Type", "application-json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{"message": "service is running"}`))
	})

	http.HandleFunc("/lightbulbs", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application-json")
		responseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(responseWriter).Encode(lightbulbs)
	})

	http.HandleFunc("/lightbulbs/switch", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application-json")

		name := request.URL.Query().Get("name")

		check := checkName(responseWriter, request, name)
		if check < 0 {
			return
		}

		lightbulbs[name] = !lightbulbs[name]

		responseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(responseWriter).Encode(lightbulbs)
	})

	http.HandleFunc("/lightbulbs/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Contet-Type", "application-json")

		name := r.URL.Query().Get("name")

		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message":"a lightbulb name should be provided as the value of a 'name' querystring parameter"}`))
			return
		}
		if _, keyExist := lightbulbs[name]; keyExist {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "a lightbulb with the provided name already exist"}`))
			return
		}

		lightbulbs[name] = false

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(lightbulbs)
	})

	http.HandleFunc("/lightbulbs/delete", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application-json")

		name := r.URL.Query().Get("name")

		check := checkName(w, r, name)
		if check < 0 {
			return
		}

		delete(lightbulbs, name)

		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(lightbulbs)

	})

	fmt.Println("http server listing on localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func checkName(responseWriter http.ResponseWriter, request *http.Request, name string) int {
	if name == "" {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{"message":"a lightbulb name should be provided as the value of a 'name' querystring parameter"}`))
		return -1
	}
	if _, keyExist := lightbulbs[name]; !keyExist {
		responseWriter.WriteHeader(http.StatusNotFound)
		responseWriter.Write([]byte(`{"message": "a lightbulb with the provided name doesn't exist"}`))
		return -2
	}
	return 1
}
