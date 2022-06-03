package match

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jarrancarr/ChexxServer/store"
	"github.com/jarrancarr/ChexxServer/user"
	"github.com/jarrancarr/ChexxServer/utils"
)

func Matches(w http.ResponseWriter, r *http.Request) {

	var m store.Match
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
		utils.Respond(w, utils.Message(false, "Error while decoding request body for new comment"))
		return
	}

	user, _ := user.FindUser(r)

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}

	match.Black.UserId = user.ID
	match.White.UserId = user.ID
	resp := match.Create()
	utils.Respond(w, resp)
}

func ListMatches(w http.ResponseWriter, r *http.Request) {

	user, _ := user.FindUser(r)
	var matches []store.Match

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}
	result := store.GetDB().Table("matches").Where("white_player_id = ? OR black_player_id = ?", user.ID, user.ID).Find(&matches)
	if result.Error != nil {
		utils.Respond(w, utils.Message(false, "Error in DB fetch for matches."))
		return
	}
	fmt.Printf("found %d matches", result.RowsAffected)

	matchesData := make([][]string, result.RowsAffected)
	for m := range matches {
		matchesData[m] = []string{"" + matches[m].Title + "", ""}
	}
	resp := utils.Message(true, "Found matches")
	resp["matches"] = matchesData

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
