package apiMiddlewares

import (
	"api/apiHandlers"
	"api/apiSructs"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

var UserInfoMiddleware = func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token := r.Context().Value("user").(*jwt.Token)
	parsedToken, err := jwt.Parse(token.Raw, func(token *jwt.Token) (interface{}, error) {
		return []byte(apiHandlers.SigningKey), nil
	})

	if err != nil {
		var answer = structs.SimpleAnswer{
			Error: true,
			Message: err.Error(),
		}
		payload, _:= json.Marshal(answer)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(payload))
		return
	}
	claims := parsedToken.Claims.(jwt.MapClaims)
	fmt.Printf("%f", claims["email"])
	//var Cred = Credential{
	//	claims["role"].(string),
	//	claims["name"].(string),
	//	claims["email"].(string),
	//}
	//fmt.Printf("%f", Cred)
	next(w, r)
}
