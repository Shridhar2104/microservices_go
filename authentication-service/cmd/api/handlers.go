package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	
)

func (app *Config) authenticate (w http.ResponseWriter, r *http.Request){
	var requestPayload struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	
	//validate the user against the database

	user, err:= app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid cred"), http.StatusBadRequest)
		return
	}
	valid, err:= app.Models.User.PasswordMatches(requestPayload.Password)
	if err != nil || !valid{
		app.errorJSON(w, errors.New("invalid cred"), http.StatusBadRequest)
		return
	}

	//log auth

	err = app.logRequest("authenticate", fmt.Sprintf("%s logged in", user.Email))

	if err!=nil{
		app.errorJSON(w, err)
		return
	}
	

	
	payload := jsonResponse{
		Success: true,
		Message: fmt.Sprintf("Logged in with uses %s", user.Email),
		Data: user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

	
}
func (app *Config) logRequest (name , data string) error{
	var entry struct{
		Name string `json:"name"`
		Data string `json:"data"`

	}

	entry.Data = data
	entry.Name = name

	jsonData, _:= json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"
	req, err:= http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-type", "application-json")
	client := &http.Client{}
	resp, err:= client.Do(req)
	if err!= nil{
		return err

	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to log request")
	}

	return nil

}