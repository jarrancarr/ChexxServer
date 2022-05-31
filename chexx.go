package main

import (
	"encoding/json"
	"fmt"
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

type Address struct {
	Street string
	City   string
}

type Person struct {
	Name    string
	Address Address
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
	w.Write([]byte(fmt.Sprintf("Gorilla!\n%s\n", r.RequestURI)))
}
func OnlineLookupHandler(w http.ResponseWriter, r *http.Request) {
	kv := r.FormValue("game")
	fmt.Println(kv)
	p := Person{
		Name: "Sherlock Holmes",
		Address: Address{
			"22/b Baker street",
			"London",
		},
	}

	res, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}
	// fmt.Println(res)

	//w.Write([]byte(fmt.Sprintf("Gorilla!\n%s\n", "hi")))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")

	//w.Write([]byte("res"))
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

var test = false

func main() {
	if test {
		store.DoTest()
		os.Exit(0)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", YourHandler)
	r.HandleFunc("/test", OnlineLookupHandler)
	r.HandleFunc("/tutor", tutor.Tutorial)
	r.HandleFunc("/user", user.UserLogin)
	r.HandleFunc("/match/getMatches", match.Matches)
	http.Handle("/", r)

	//log.Fatal(http.ListenAndServe(":8000", r))
	log.Fatal(http.ListenAndServe(":8000",
		handlers.LoggingHandler(os.Stdout, handlers.CORS(
			handlers.AllowedMethods([]string{"POST", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "X-Requested-With"}))(r))))
}
