package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"strings"
	"time"
)

type SimpleAnswer struct {
	Error bool `json:"error"`
	Message string `json:"message"`
}

var NotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	var answer = SimpleAnswer{
		Error:   true,
		Message: "Method's not allowed",
	}
	payload, _ := json.Marshal(answer)
	w.Write([]byte(payload))
})

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	var answer = SimpleAnswer{
		Error:   true,
		Message: r.Host,
	}
	payload, _ := json.Marshal(answer)
	w.Write([]byte(payload))
})

var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	var answer = SimpleAnswer{
		Error: false,
		Message: "Engineer Api's up and running",
	}
	payload, _:= json.Marshal(answer)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

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

var checkCredentials = func(login, password string) (Credential, error){
	//byt := []byte(`{"role":"admin","name":"Вася Пупкинович"}`)
	byt := []byte(`{"role":"admin","name":"Вася Пупкинович", "email":"nvs@nvs.ru"}`)
	res := Credential{}
	json.Unmarshal(byt, &res)
	var err error
	err = nil

	if len(res.Email) <= 0 {
		err = errors.New(fmt.Sprintf("BingoBoom auth server doesn't send required field 'email' for user %d. Sorry. Can't authenticate.", login))
	}

	//err = "login and password are incorrect";

	return res, err
}

var signingKey = []byte("Engineer_Panel_Secret_Key_159753")

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	user, err := checkCredentials("ads", "asdasdads")
	if err != nil {
		var answer = SimpleAnswer{
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
	tokenString, _ := token.SignedString(signingKey)

	var answer = AuthDataAnswer{
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

var userModelMiddleware = func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token := r.Context().Value("user").(*jwt.Token)
	parsedToken, err := jwt.Parse(token.Raw, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	if err != nil {
		var answer = SimpleAnswer{
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

var jwtApiMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error){
		return []byte(signingKey), nil
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
		var answer = SimpleAnswer{
			Error: true,
			Message: err,
		}
		payload, _:= json.Marshal(answer)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(payload))
	},
})

func main(){
	/* Основной Роут для простых url*/
	r:= mux.NewRouter()
	r.Handle("/status", StatusHandler).Methods("GET")
	r.Handle("/login", GetTokenHandler).Methods("POST")

	/* Суброутер для всех url, которые начинается с /api */
	api := mux.NewRouter().PathPrefix("/api").Subrouter()

	/*
	Суброутер для всех url, которые начинается с /api
	будет защищено c помощью jwt middleware
	*/
	r.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(jwtApiMiddleware.HandlerWithNext),
		negroni.HandlerFunc(userModelMiddleware),
		negroni.Wrap(api),
	))

	api.HandleFunc("/test", testHandler)

	n := negroni.Classic()
	n.UseHandler(r)

	r.PathPrefix("/").HandlerFunc(NotAllowed)

	/* CORS Разрешаем заголовки */
	headers := handlers.AllowedHeaders([]string{"Content-Type", "X-Auth-Token"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	http.ListenAndServe(":5000", handlers.CORS(headers,methods, origins)(n))
}