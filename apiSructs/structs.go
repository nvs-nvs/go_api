package structs

import "github.com/gorilla/mux"

type Server struct {
	//db     *someDatabase
	router *mux.Router
}

type SimpleAnswer struct {
	Error bool `json:"error"`
	Message string `json:"message"`
}

type Credential struct {
	Role string `json:"role"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type AuthDataAnswer struct {
	*Credential `json:"user"`
	Token string `json:"token"`
	Error bool `json:"error"`
	Expires string `json:"expires"`
	Auth bool `json:"auth"`
}