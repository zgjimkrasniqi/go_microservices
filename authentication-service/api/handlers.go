package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// A simple way to decode request body
	// err := json.NewDecoder(r.Body).Decode(&requestPayload)

	// Decoding request body using helpers
	err := app.readJson(r, &requestPayload)

	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		_ = app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		_ = app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, _ := app.Models.User.GetAllUsers()
	log.Println("users in the database", users)

	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Users in the database"),
		Data:    users,
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}
