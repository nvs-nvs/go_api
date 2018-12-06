package apiHandlers

import (
	"api/apiSructs"
	"encoding/json"
	"net/http"
)

var TestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	var answer = structs.SimpleAnswer{
		Error:   true,
		Message: r.Host,
	}
	payload, _ := json.Marshal(answer)
	w.Write([]byte(payload))
})

var NotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	var answer = structs.SimpleAnswer{
		Error:   true,
		Message: "Method's not allowed",
	}
	payload, _ := json.Marshal(answer)
	w.Write([]byte(payload))
})

var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var answer = structs.SimpleAnswer{
		Error:   false,
		Message: "Engineer Api's up and running",
	}
	payload, _ := json.Marshal(answer)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})