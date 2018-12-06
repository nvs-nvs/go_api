package apiMiddlewares

import (
	"api/apiHandlers"
	"api/apiSructs"
	"encoding/json"
	"errors"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

var JwtApiMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error){
		return []byte(apiHandlers.SigningKey), nil
	},
	Extractor: func(r *http.Request) (s string, e error) {
		authHeader := r.Header.Get("X-Auth-Token")
		if authHeader == "" {
			return "", errors.New("Token is required")
		}
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			return "", errors.New("Authorization header format must be Bearer {token}")
		}
		return authHeaderParts[1], nil
	},
	SigningMethod: jwt.SigningMethodHS256,
	ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
		var answer = structs.SimpleAnswer{
			Error: true,
			Message: err,
		}
		payload, _:= json.Marshal(answer)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(payload))
	},
})