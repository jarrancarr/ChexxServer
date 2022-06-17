package main

import (
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
	"github.com/jarrancarr/ChexxServer/utils"
)

var test = false

func main() {
	utils.Init()
	if test {
		fmt.Println("running test")
		//store.DoTest()
		match := &store.Match{}
		//match.White.Pieces = []string{"Kc44", "Qd41", "Id31", "Ed4", "Rd54", "Rd5", "Rc52", "Nd53", "Nd51", "Nc33", "Bc53", "Bc55", "Bd52", "Ad42", "Ad3", "Ac43", "Pd55", "Pd44", "Pd33", "Pd21", "Pc22", "Pc31", "Pc41", "Pc51", "Sd43", "Sd32", "Sd2", "Sc32", "Sc42"}
		//match.Black.Pieces = []string{"Ka41", "Qf44", "If33", "Ea4", "Ra5", "Rf52", "Ra54", "Nf53", "Nf55", "Na31", "Ba53", "Ba51", "Bf54", "Af43", "Aa3", "Aa42", "Pf51", "Pf41", "Pf31", "Pf22", "Pa21", "Pa33", "Pa44", "Pa55", "Sf42", "Sf32", "Sa2", "Sa32", "Sa43"}

		//                             Nd53    Nd51    Bc53    Bc55    Bd52    Kc54    Pd55    Pd44    Pe21    Pe    P*    Pc31    Pc41    Pc51    Sd43    Se1    Sc1    Sc42    Ad42    Ae2    Ac4
		// 							   Ka52    If2    Ec2    Pf51    Pf41    Pe22    Pf    Pa33    Pa44    Pa55    Sf42    Sf1    Sa    Sb1    Sa43    Af43    Aa3    Aa42

		match.White.Pieces = []string{"Bc53", "Bc55", "Kc54", "Ac43"}
		match.Black.Pieces = []string{"Ka52", "Ea11"}
		//match.White.Pieces = []string{"Kc44", "Qd41", "Sd32"}
		//match.Black.Pieces = []string{"Ka41", "Pf21", "Bc41"}
		match.Log = []string{"blank"} // to make this blacks move
		// match.Move("a2~a1")
		// match.TestAttacks("f5")
		// match.Show("f5")
		// match.Analyse()
		// match.Examine()
		for i := 0; i < 1; i++ {
			best := match.AI(4, 2, nil)
			match.Move(best.LastMove)
		}
		fmt.Println("move: " + match.LastMove)
		// match.Analyse()
		//match.Examine()
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
	r.HandleFunc("/user/{id}", user.UserInfo)
	r.HandleFunc("/match/getMatches", match.Matches)
	r.HandleFunc("/match/save", match.SaveMatch)
	r.HandleFunc("/match/list", match.ListMatches)
	r.HandleFunc("/match/load/{id}", match.LoadMatch)
	r.HandleFunc("/match/move/{id}", match.MakeMove)
	r.HandleFunc("/match/challenge", match.CreateMatch)
	r.HandleFunc("/match/accept/{id}", match.AcceptMatch)
	r.HandleFunc("/match/delete/{id}", match.DeleteMatch)
	r.HandleFunc("/match/cpu/{level}", match.AIMove)

	http.Handle("/", r)
	r.Use(user.JwtAuthentication)

	router := handlers.LoggingHandler(os.Stdout, handlers.CORS(
		handlers.AllowedMethods([]string{"POST", "OPTIONS", "GET"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Authorization", "Content-Type", "Content-Language", "Origin", "X-Requested-With"}))(r))

	log.Fatal(http.ListenAndServe(":8000", router))

	// log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", router))

	// mgr := autocert.Manager{
	// 	Cache:      autocert.DirCache("certCache"),
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist("chexx.org", "www.chexx.org"),
	// }

	// server := &http.Server{
	// 	Addr:      ":https",
	// 	Handler:   router,
	// 	TLSConfig: mgr.TLSConfig(),
	// }

	// server.ListenAndServeTLS("", "")
}
