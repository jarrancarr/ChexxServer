package match

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jarrancarr/ChexxServer/store"
	"github.com/jarrancarr/ChexxServer/utils"
)

func Matches(w http.ResponseWriter, r *http.Request) {

	var m store.MatchWrap
	err := json.NewDecoder(r.Body).Decode(&m)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}
	fmt.Println("------------------------------")
	fmt.Println("------------------------------")
	utils.CorsHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("good"))
}

func SaveMatch(w http.ResponseWriter, r *http.Request) {
	match := &store.Match{}
	err := json.NewDecoder(r.Body).Decode(match)
	if err != nil {
		fmt.Println(err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body for new comment"))
		return
	}

	fmt.Println("---------------SaveMatch--------------")
	fmt.Printf("%v\n", match)
	fmt.Println("---------------SaveMatch--------------")

	resp := match.Create()
	utils.Respond(w, resp)
}

// var StoreComment = func(w http.ResponseWriter, r *http.Request) {
// 	user, err := user.FindUser(r)

// 	com := &models.Comment{}
// 	err = json.NewDecoder(r.Body).Decode(com)
// 	if err != nil {
// 		fmt.Println(err)
// 		u.Respond(w, u.Message(false, "Error while decoding request body for new comment"))
// 		return
// 	}
// 	com.Author = user.ID
// 	com.RefType = models.GENERAL
// 	//user, _ := models.FindUser(r)
// 	resp := com.Create()
// 	u.Respond(w, resp)
// }
