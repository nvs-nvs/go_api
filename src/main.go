package src

import (
	"encoding/json"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/urfave/negroni"
	"net/http"
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

var api = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
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
}

var checkCredentials = func(login, password string) (cred Credential, err error){
	return Credential{
		Role: "admin",
		Name: "Вася Пупкин",
	},
	errors.New("gggggggggg")
}

var signingKey = []byte("Engineer_Panel_Secret_Key_159753")

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	user, err := checkCredentials("ads", "asdasdads")
	if err != nil {
		var answer = SimpleAnswer{
			Error: false,
			Message: "login and password are incorrect",
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
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token.Claims = claims
	tokenString, _ := token.SignedString(signingKey)
	w.Write([]byte(tokenString))
})

var jwtApiMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error){
		return signingKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
	ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
		var answer = SimpleAnswer{
			Error: true,
			Message: "Token is required but not exists",
		}
		payload, _:= json.Marshal(answer)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(payload))
	},
})

func main(){
	/* Простые роуты */
	r:= mux.NewRouter()
	r.Handle("/status", StatusHandler).Methods("GET")
	r.Handle("/login", GetTokenHandler).Methods("POST", "GET")

	/* Роуты Апи */
	apiRouter := mux.NewRouter()

	apiRouter.Handle("/test", api).Methods("GET", "POST")

	negroniInstance :=
		negroni.New(
			negroni.HandlerFunc(
				jwtApiMiddleware.HandlerWithNext),
				negroni.Wrap(apiRouter))

	r.PathPrefix("/api").Handler(negroniInstance)

	r.PathPrefix("/").HandlerFunc(NotAllowed)

	/* CORS Разрешаем заголовки */
	headers := handlers.AllowedHeaders([]string{"Content-Type", "X-Auth-Token"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	http.ListenAndServe(":5000", handlers.CORS(headers,methods, origins)(r))
}