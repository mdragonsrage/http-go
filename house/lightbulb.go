package house

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var badRequestResponse = []byte(`{"message": "bad request"}`)
var methodNotAllowedResponse = []byte(`{"message": "method not allowed"}`)
var notFoundResponse = []byte(`{"message": "lightbulb not found"}`)

type Storage interface {
	GetAll() ([]Lightbulb, error)
	Get(name string) (Lightbulb, error)
	Create(lb Lightbulb) error
	Update(lb Lightbulb) error
	Delete(name string) error
}

type Lightbulb struct {
	Name string `json:"name"`
	On   bool   `json:"on"`
}

func GetLightbulb(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application-json")
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(methodNotAllowedResponse)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			lightbulbs, err := storage.GetAll()
			if err != nil {
				msg := fmt.Sprintf("an error occurred while trying to getall lightbulbs: %v\n", err)
				log.Println(msg)
				w.WriteHeader((http.StatusInternalServerError))
				w.Write([]byte(fmt.Sprintf(`{"message": %s}`, msg)))
				return
			}

			response := map[string]bool{}
			for _, lightbulb := range lightbulbs {
				response[lightbulb.Name] = lightbulb.On
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		lightbulb, err := storage.Get(name)
		if err != nil {
			msg := fmt.Sprintf("an error occurred while trying to get lightbulb: %v\n", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, msg)))
			return
		}

		json.NewEncoder(w).Encode(lightbulb)
	}
}

func CreateLightbulb(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application-json")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(methodNotAllowedResponse)
			return
		}

		if r.Body == nil {
			log.Println("create requires a request body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}
		var lightbulb Lightbulb
		err := json.NewDecoder(r.Body).Decode(&lightbulb)
		if err != nil {
			msg := fmt.Sprintf("an error occurred while trying to create lightbulb: %v\n", err)
			log.Println(msg)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		err = storage.Create(lightbulb)
		if err != nil {
			msg := fmt.Sprintf("an error occurred while trying to create lightbulbs: %v\n", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, msg)))
			return
		}

		json.NewEncoder(w).Encode(lightbulb)
	}
}

func SwitchLightbulb(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(methodNotAllowedResponse)
			return
		}

		name := r.URL.Query().Get("name")
		lightbulb, err := storage.Get(name)
		if err != nil {
			msg := fmt.Sprintf("an error occurred while trying to get lightbulb: %v\n", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, msg)))
			return
		}

		lightbulb.On = !lightbulb.On

		err = storage.Update(lightbulb)
		if err != nil {
			msg := fmt.Sprintf("an error occurred while trying to update lightbulbs: %v\n", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, msg)))
			return
		}

		lightbulbs, err := storage.GetAll()
		if err != nil {
			msg := fmt.Sprintf("an error occurred while trying to getall lightbulbs: %v\n", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, msg)))
			return
		}

		json.NewEncoder(w).Encode(lightbulbs)
	}
}

func DeleteLightbulb(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(methodNotAllowedResponse)
			return
		}

		name := r.URL.Query().Get("name")

		err := storage.Delete(name)
		if err != nil {
			msg := fmt.Sprintf("an error occurred while trying to delete lightbulbs: %v\n", err)
			log.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, msg)))
			return
		}
	}
}
