package match

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

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
	utils.CorsHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("good"))
}
func DeleteMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Match ID not found."))
		return
	}
	m := store.GetMatch(uint(id))

	user, _ := user.FindUser(r)
	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}
	store.GetDB().Delete(m)

	ListMatches(w, r)
}

func AIMove(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	dat := params["level"]
	level, err := strconv.Atoi(dat)
	if err != nil {
		utils.Respond(w, utils.Message(false, "No level found"))
		level = 1
	}

	match := &store.Match{}
	err = json.NewDecoder(r.Body).Decode(match)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body for new comment"))
		return
	}

	user, _ := user.FindUser(r)

	if user == nil {
		if level > 3 {
			level = 3
		}
		move := match.AI(9, level, nil)

		resp := utils.Message(true, "Checkmate or something.")
		resp["move"] = move.LastMove
		utils.Respond(w, resp)
		//utils.Respond(w, utils.Message(false, "User not found."))
		return
	}

	move := match.AI(9, level, nil)

	resp := utils.Message(true, "Checkmate or something.")
	resp["move"] = move.LastMove
	utils.Respond(w, resp)
}

func AcceptMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Match ID not found."))
		return
	}
	m := store.GetMatch(uint(id))

	user, _ := user.FindUser(r)
	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}

	if m.BlackPlayerId == 0 {
		m.BlackPlayerId = user.ID
	} else {
		m.WhitePlayerId = user.ID
	}

	m.Game.Status = "engaged"
	store.GetDB().Save(m)

	ListMatches(w, r)
}

func CreateMatch(w http.ResponseWriter, r *http.Request) {

	var jsonData map[string]string
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body for new challenge"))
		return
	}
	user, _ := user.FindUser(r)

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}

	match := &store.Match{}

	if jsonData["color"] == "0" {
		match.Black.UserId = user.ID
	}
	if jsonData["color"] == "1" {
		match.White.UserId = user.ID
	}
	if jsonData["color"] == "2" {
		coin := rand.Intn(2)
		if coin == 0 {
			match.Black.UserId = user.ID
		} else {
			match.White.UserId = user.ID
		}
	}
	switch jsonData["dpm"] {
	case "2":
		match.Game = store.Type{Name: "2 day", GameClock: 3600 * 48, MoveClock: 0, Status: "open", Rating: user.Rank}
		break
	case "3":
		match.Game = store.Type{Name: "3 day", GameClock: 3600 * 72, MoveClock: 0, Status: "open", Rating: user.Rank}
		break
	case "5":
		match.Game = store.Type{Name: "5 day", GameClock: 3600 * 120, MoveClock: 0, Status: "open", Rating: user.Rank}
		break
	case "7":
		match.Game = store.Type{Name: "7 day", GameClock: 3600 * 168, MoveClock: 0, Status: "open", Rating: user.Rank}
		break
	case "10":
		match.Game = store.Type{Name: "10 day", GameClock: 3600 * 240, MoveClock: 0, Status: "open", Rating: user.Rank}
		break
	case "14":
		match.Game = store.Type{Name: "14 day", GameClock: 3600 * 336, MoveClock: 0, Status: "open", Rating: user.Rank}
		break
	}
	match.White.Pieces = []string{"Rd54", "Rd5", "Rc52", "Nd53", "Nd51", "Nc33", "Bc53", "Bc55", "Bd52", "Qd41", "Kc44", "Id31", "Ed4", "Pd55", "Pd44", "Pd33", "Pd21", "Pc22", "Pc31", "Pc41", "Pc51", "Sd43", "Sd32", "Sd2", "Sc32", "Sc42", "Ad42", "Ad3", "Ac43"}
	match.Black.Pieces = []string{"Ra5", "Rf52", "Ra54", "Nf53", "Nf55", "Na31", "Ba53", "Ba51", "Bf54", "Qf44", "Ka41", "If33", "Ea4", "Pf51", "Pf41", "Pf31", "Pf22", "Pa21", "Pa33", "Pa44", "Pa55", "Sf42", "Sf32", "Sa2", "Sa32", "Sa43", "Af43", "Aa3", "Aa42"}
	match.Title = jsonData["title"]
	match.Create()
	ListMatches(w, r)
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
	if match.ID == 0 {
		resp := match.Create()
		utils.Respond(w, resp)
		return
	}
	resp := match.Update()
	utils.Respond(w, resp)
}

func MakeMove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Match ID not found."))
		return
	}
	m := store.GetMatch(uint(id))

	user, _ := user.FindUser(r)

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}

	// is it your move?
	if m.Logs == "" || len(strings.Split(strings.Trim(m.Logs, " "), " "))%2 == 0 { // white turn
		if m.WhitePlayerId != user.ID {
			utils.Respond(w, utils.Message(false, "Not your turn, dumbass."))
			return
		}
	} else {
		if m.BlackPlayerId != user.ID {
			utils.Respond(w, utils.Message(false, "Not your turn, dumbass."))
			return
		}
	}
	var jsonData map[string]string
	err = json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body for move"))
		return
	}
	move := jsonData["move"]
	if move == "" {
		utils.Respond(w, utils.Message(false, "cannot get move from json"))
		return
	}
	m.Move(move)
	resp := m.Update()
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

	if m.BlackPlayerId != 0 {
		resp["black"] = store.GetUser(m.BlackPlayerId)
	}
	if m.WhitePlayerId != 0 {
		resp["white"] = store.GetUser(m.WhitePlayerId)
	}

	utils.Respond(w, resp)
}

func ListMatches(w http.ResponseWriter, r *http.Request) {

	user, _ := user.FindUser(r)
	var matches []store.Match
	var open []store.Match

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}
	store.GetDB().Table("matches").Where("white_player_id = ? OR black_player_id = ?", user.ID, user.ID).Find(&matches)
	store.GetDB().Table("matches").Where("white_player_id != ? AND black_player_id != ? AND status='open' AND (rating>=? AND rating <=?) ", user.ID, user.ID, user.Rank-200, user.Rank+200).Find(&open)

	savedMatches := make([][]string, 0)
	myOpenMatches := make([][]string, 0)
	openMatches := make([][]string, 0)
	readyMatches := make([][]string, 0)
	waitingMatches := make([][]string, 0)
	for m := range matches {
		if matches[m].BlackPlayerId == user.ID && matches[m].WhitePlayerId == user.ID {
			savedMatches = append(savedMatches, []string{fmt.Sprintf(".%s.:%d", matches[m].Title, matches[m].ID), ""})
		} else if matches[m].BlackPlayerId == user.ID || matches[m].WhitePlayerId == user.ID { // so only one player is user
			if matches[m].Game.Status == "open" {
				myOpenMatches = append(myOpenMatches, []string{fmt.Sprintf(".%s.:%d", matches[m].Title, matches[m].ID), ""})
			} else if matches[m].Game.Status == "engaged" {
				if matches[m].Logs == "" || len(strings.Split(strings.Trim(matches[m].Logs, " "), " "))%2 == 0 { // white turn
					if matches[m].WhitePlayerId == user.ID {
						readyMatches = append(readyMatches, []string{fmt.Sprintf(".%s.:%d", matches[m].Title, matches[m].ID), ""})
					} else {
						waitingMatches = append(waitingMatches, []string{fmt.Sprintf(".%s.:%d", matches[m].Title, matches[m].ID), ""})
					}
				} else { // black turn
					if matches[m].BlackPlayerId == user.ID {
						readyMatches = append(readyMatches, []string{fmt.Sprintf(".%s.:%d", matches[m].Title, matches[m].ID), ""})
					} else {
						waitingMatches = append(waitingMatches, []string{fmt.Sprintf(".%s.:%d", matches[m].Title, matches[m].ID), ""})
					}
				}
			}
		}
	}
	for m := range open {
		openMatches = append(openMatches, []string{fmt.Sprintf(".%s.:%d:%d", open[m].Title, open[m].ID, open[m].Game.Rating), ""})
	}
	resp := utils.Message(true, "Found matches")
	resp["savedMatches"] = savedMatches
	resp["myOpen"] = myOpenMatches
	resp["open"] = openMatches
	resp["ready"] = readyMatches
	resp["waiting"] = waitingMatches
	// resp["victory"] = openMatches
	// resp["defeat"] = openMatches

	utils.Respond(w, resp)
}
