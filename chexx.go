package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jarrancarr/ChexxServer/match"
	"github.com/jarrancarr/ChexxServer/store"
	"github.com/jarrancarr/ChexxServer/tutor"
	"github.com/jarrancarr/ChexxServer/user"
)

var test = false

func main() {
	if test {
		store.DoTest()
		os.Exit(0)
	}
	r := mux.NewRouter()
	// c := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	AllowCredentials: true,
	// 	AllowedHeaders: []string{"Authorization"},
	// 	// Enable Debugging for testing, consider disabling in production
	// 	Debug: true,
	// })
	r.HandleFunc("/tutor", tutor.Tutorial)
	r.HandleFunc("/user/login", user.Authenticate)
	r.HandleFunc("/user/facebook", user.Facebook)
	r.HandleFunc("/user/logout", user.Logout)
	r.HandleFunc("/user/register", user.RegisterUser)
	r.HandleFunc("/match/getMatches", match.Matches)
	r.HandleFunc("/match/save", match.SaveMatch)
	r.HandleFunc("/match/list", match.ListMatches)
	http.Handle("/", r)

	r.Use(user.JwtAuthentication)

	//log.Fatal(http.ListenAndServe(":8000", r))

	log.Fatal(http.ListenAndServe(":8001",
		handlers.LoggingHandler(os.Stdout, handlers.CORS(
			handlers.AllowedMethods([]string{"POST", "OPTIONS", "GET"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Authorization", "Content-Type", "Content-Language", "Origin", "X-Requested-With"}))(r))))
}
