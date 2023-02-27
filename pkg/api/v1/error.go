package v1

import (
	"encoding/json"
	"log"
	"net/http"
)

// Errors ...
func Errors(w http.ResponseWriter, statusCode int, errors []string) {
	w.WriteHeader(statusCode)

	if errors == nil {
		write(w, nil)
		return
	}

	p := map[string][]string{
		"error": errors,
	}
	data, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}

	write(w, data)
}

// Error ...
func Error(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)

	var p map[string]string
	if message == "" {
		write(w, nil)
		return
	}

	p = map[string]string{
		"error": message,
	}
	data, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	write(w, data)
}

func write(w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
