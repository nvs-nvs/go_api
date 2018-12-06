package main

import (
	"api/apiHandlers"
	"api/apiMiddlewares"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)



func main() {
	/* Основной Роут для простых url*/
	r := mux.NewRouter()
	r.Handle("/status", apiHandlers.StatusHandler).Methods("GET")
	r.Handle("/login", apiHandlers.GetTokenHandler).Methods("POST")

	/* Суброутер для всех url, которые начинается с /api */
	api := mux.NewRouter().PathPrefix("/api").Subrouter()

	/*
	Суброутер для всех url, которые начинается с /api
	будет защищено c помощью jwt middleware
	*/
	r.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(apiMiddlewares.JwtApiMiddleware.HandlerWithNext),
		negroni.HandlerFunc(apiMiddlewares.UserInfoMiddleware),
		negroni.Wrap(api),
	))

	api.HandleFunc("/test", apiHandlers.TestHandler)

	n := negroni.Classic()

	n.UseHandler(r)

	r.PathPrefix("/").HandlerFunc(apiHandlers.NotAllowed)

	/* CORS Разрешаем заголовки */
	headers := handlers.AllowedHeaders([]string{"Content-Type", "X-Auth-Token"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	http.ListenAndServe(":5000", handlers.CORS(headers, methods, origins)(n))
}
