package main

import (
	"fmt"
	"logger-service/data"
	"net/http"
)

func (app *Config) Log(w http.ResponseWriter, r *http.Request) {
	var requestPayload data.LogEntry

	// A simple way to decode request body
	// err := json.NewDecoder(r.Body).Decode(&requestPayload)

	// Decoding request body using helpers
	err := app.readJson(r, &requestPayload)

	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log, err := app.LogEntry.Insert(requestPayload)

	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Log Id: %s", log),
		Data:    log,
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetAllLogs(w http.ResponseWriter, r *http.Request) {
	logs, _ := app.LogEntry.GetAllLogs()

	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logs in the database"),
		Data:    logs,
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}
