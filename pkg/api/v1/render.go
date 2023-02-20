package v1

import (
	"encoding/json"
	"log"
	"net/http"
)

// ResponseBase ...
type ResponseBase struct {
	Data interface{} `json:"data"`
	Meta MetaData    `json:"meta,omitempty"`
}

// MetaData ...
type MetaData struct {
	Size  int `json:"size"`
	Total int `json:"total"`
}

// SendJSON ...
func SendJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.WriteHeader(statusCode)

	if payload == nil {
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		WithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		WithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

// Status ...
func Status(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}
