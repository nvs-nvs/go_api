package apiHandlers

import (
	"api/apiSructs"
	"api/apiUtils"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var SigningKey = []byte("Engineer_Panel_Secret_Key_159753")

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	user, err := apiUtils.CheckCredentials("ads", "asdasdads")
	if err != nil {
		var answer = structs.SimpleAnswer{
			Error: false,
			Message: err.Error(),
		}
		payload, _:= json.Marshal(answer)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(payload))
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)

	claims["role"] = user.Role
	claims["name"] = user.Name
	expires := time.Now().Add(time.Hour*24)

	claims["exp"] = expires.Unix()
	token.Claims = claims
	tokenString, _ := token.SignedString(SigningKey)

	var answer = structs.AuthDataAnswer{
		&user,
		tokenString,
		false,
		expires.Format("2006-01-02 15:04:05"),
		true,
	}
	payload, _:= json.Marshal(answer)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})
