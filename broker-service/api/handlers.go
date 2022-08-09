package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := JsonResponse{Error: false, Message: "Hit the broker"}

	/*	response, _ := json.Marshal(payload)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write(response)
	*/

	// Using helpers
	_ = app.writeJSON(w, http.StatusOK, payload)
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJson(r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		{
			data, _ := json.Marshal(requestPayload)
			l := LogPayload{
				Name: requestPayload.Action,
				Data: string(data),
			}
			app.authenticate(w, requestPayload.Auth)
			app.insertLog(w, l)
		}

	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// Create some json and send to the auth microservice
	jsonData, _ := json.Marshal(a)

	// Call the service
	// url: http://<name of the service specified in docker-compose.yml>/...
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))

	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		_ = app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		_ = app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// Create a variable to read response.Body into
	var jsonFromService JsonResponse

	// Decode the Json
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		_ = app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload JsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) insertLog(w http.ResponseWriter, l LogPayload) {
	// Create some json and send to the logger microservice
	jsonData, _ := json.Marshal(l)

	// Call the service
	// url: http://<name of the service specified in docker-compose.yml>/...
	request, err := http.NewRequest("POST", "http://logger-service/insert_log", bytes.NewBuffer(jsonData))

	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()
}
