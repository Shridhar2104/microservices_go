package main

import (
	"logger-service/data"
	"net/http"

)

type JSONPayload struct{
	Name string `json:"name"`
	Data string `json:"data"`
}
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request){
	var requestPayLoad JSONPayload

	_ = app.readJSON(w, r, &requestPayLoad)
	event:= data.LogEntry{
		Name: requestPayLoad.Name,
		Data: requestPayLoad.Data,
	}
	err:= app.Models.LogEntry.Insert(event)

	if err!=nil{
		app.errorJSON(w, err)
	}

	resp:= jsonResponse{
		Success: true,
		Message: "Log Entry Created",
	}
	app.writeJSON(w, http.StatusAccepted,resp)
}