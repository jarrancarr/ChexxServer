package match

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func LoadMatch(w http.ResponseWriter, r *http.Request) { // {id:0, name:'offline', white:{pieces:['Rd54', 'Rd5', 'Rc52', 'Nd53', 'Nd51', 'Nc33', 'Bc53', 'Bc55', 'Bd52', 'Qd41', 'Kc44', 'Id31', 'Ed4', 'Pd55', 'Pd44', 'Pd33', 'Pd21', 'Pc22', 'Pc31', 'Pc41', 'Pc51', 'Sd43', 'Sd32', 'Sd2', 'Sc32', 'Sc42', 'Ad42', 'Ad3', 'Ac43'], time:300}, black:{pieces:['Ra5', 'Rf52', 'Ra54', 'Nf53', 'Nf55', 'Na31', 'Ba53', 'Ba51', 'Bf54', 'Qf44', 'Ka41', 'If33', 'Ea4', 'Pf51', 'Pf41', 'Pf31', 'Pf22', 'Pa21', 'Pa33', 'Pa44', 'Pa55', 'Sf42', 'Sf32', 'Sa2', 'Sa32', 'Sa43', 'Af43', 'Aa3', 'Aa42'], time:300}, log:[], type:{game:300, move:15}});
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Match ID not found."))
		return
	}
	m := store.GetMatch(uint(id))

	resp := utils.Message(true, "Found match")
	resp["match"] = m

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
		matchesData[m] = []string{fmt.Sprintf("%s:%d", matches[m].Title, matches[m].ID), ""}
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
