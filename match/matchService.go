package match

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jarrancarr/ChexxServer/store"
	"github.com/jarrancarr/ChexxServer/user"
	"github.com/jarrancarr/ChexxServer/utils"
)

var DEBUG = false

func Matches(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("Matches")
	}
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
	if DEBUG {
		log.Println("DeleteMatch")
	}
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
	if DEBUG {
		log.Println("AIMove")
	}

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

		resp := utils.Message(true, "OK")
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
	if DEBUG {
		log.Println("AcceptMatch")
	}
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
	if DEBUG {
		log.Println("CreateMatch")
	}

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
		match.Game = store.Type{Name: "2-Day/Move", GameClock: 3600 * 48, MoveClock: 0, Status: "open", Rating: user.Rating("2-Day/Move")}
		break
	case "3":
		match.Game = store.Type{Name: "3-Day/Move", GameClock: 3600 * 72, MoveClock: 0, Status: "open", Rating: user.Rating("3-Day/Move")}
		break
	case "5":
		match.Game = store.Type{Name: "5-Day/Move", GameClock: 3600 * 120, MoveClock: 0, Status: "open", Rating: user.Rating("5-Day/Move")}
		break
	case "7":
		match.Game = store.Type{Name: "7-Day/Move", GameClock: 3600 * 168, MoveClock: 0, Status: "open", Rating: user.Rating("7-Day/Move")}
		break
	case "10":
		match.Game = store.Type{Name: "10-Day/Move", GameClock: 3600 * 240, MoveClock: 0, Status: "open", Rating: user.Rating("10-Day/Move")}
		break
	case "14":
		match.Game = store.Type{Name: "14-Day/Move", GameClock: 3600 * 336, MoveClock: 0, Status: "open", Rating: user.Rating("14-Day/Move")}
		break
	}
	match.White.Pieces = []string{"Rd54", "Rd5", "Rc52", "Nd53", "Nd51", "Nc33", "Bc53", "Bc55", "Bd52", "Qd41", "Kc44", "Id31", "Ed4", "Pd55", "Pd44", "Pd33", "Pd21", "Pc22", "Pc31", "Pc41", "Pc51", "Sd43", "Sd32", "Sd2", "Sc32", "Sc42", "Ad42", "Ad3", "Ac43"}
	match.Black.Pieces = []string{"Ra5", "Rf52", "Ra54", "Nf53", "Nf55", "Na31", "Ba53", "Ba51", "Bf54", "Qf44", "Ka41", "If33", "Ea4", "Pf51", "Pf41", "Pf31", "Pf22", "Pa21", "Pa33", "Pa44", "Pa55", "Sf42", "Sf32", "Sa2", "Sa32", "Sa43", "Af43", "Aa3", "Aa42"}
	match.Title = jsonData["title"]
	match.Create()
	ListMatches(w, r)
}
func ResignMatch(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("ResignMatch")
	}
	var mId uint
	err := json.NewDecoder(r.Body).Decode(&mId)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body for new comment"))
		return
	}

	m := store.GetMatch(mId)
	user, _ := user.FindUser(r)

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}

	if m.Logs == "" || len(strings.Split(strings.Trim(m.Logs, " "), ":::"))%2 == 0 { // white turn
		m.WinLoseDraw("lose", "White Resigns")
	} else {
		m.WinLoseDraw("win", "Black Resigns")
	}

	resp := m.Update()

	resp["rank"] = fmt.Sprintf("Your new rating is %d", user.Rating(m.Game.Name))

	utils.Respond(w, resp)
}
func DrawMatch(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("DrawMatch")
	}
	var mId uint
	err := json.NewDecoder(r.Body).Decode(&mId)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body for new comment"))
		return
	}

	m := store.GetMatch(mId)
	user, _ := user.FindUser(r)

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}

	if m.Game.Status == "Draw Offered" {
		m.WinLoseDraw("draw", "Draw")
	} else {
		m.Game.Status = "Draw Offered"
	}

	resp := m.Update()

	resp["rank"] = fmt.Sprintf("Your new rating is %d", user.Rating(m.Game.Name))

	utils.Respond(w, resp)
}
func SaveMatch(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("SaveMatch")
	}
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
	if DEBUG {
		log.Println("MakeMove")
	}
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Match ID not found."))
		return
	}
	m := store.GetMatch(uint(id))

	user, _ := user.FindUser(r)
	opponentToken := store.Online()[m.BlackPlayerId]
	if user.ID == m.BlackPlayerId {
		opponentToken = store.Online()[m.WhitePlayerId]
	}

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}

	// is it your move?
	if m.Logs == "" || len(strings.Split(strings.Trim(m.Logs, " "), ":::"))%2 == 0 { // white turn
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
	m.Move(move, true)
	resp := m.Update()
	//resp["state"] = m.AI(-1, 1, nil).LastMove
	// if m.LastMove == "Checkmate" {
	// 	lastLog := m.Log[len(m.Log)-1]
	// 	m.Log = append(m.Log[:len(m.Log)-2], lastLog+"++")
	// 	m.Update()
	// } else if m.LastMove == "Stalemate" {
	// 	lastLog := m.Log[len(m.Log)-1]
	// 	m.Log = append(m.Log[:len(m.Log)-2], lastLog+"=")
	// 	m.Update()
	// }
	if store.SessionMap[opponentToken] != nil {
		store.SessionMap[opponentToken].Inbox <- fmt.Sprintf("type||notify|||match||%s-%d|||state||"+m.LastMove, m.Title, m.ID)
	}
	utils.Respond(w, resp)
}
func LoadMatch(w http.ResponseWriter, r *http.Request) { // {id:0, name:'offline', white:{pieces:['Rd54', 'Rd5', 'Rc52', 'Nd53', 'Nd51', 'Nc33', 'Bc53', 'Bc55', 'Bd52', 'Qd41', 'Kc44', 'Id31', 'Ed4', 'Pd55', 'Pd44', 'Pd33', 'Pd21', 'Pc22', 'Pc31', 'Pc41', 'Pc51', 'Sd43', 'Sd32', 'Sd2', 'Sc32', 'Sc42', 'Ad42', 'Ad3', 'Ac43'], time:300}, black:{pieces:['Ra5', 'Rf52', 'Ra54', 'Nf53', 'Nf55', 'Na31', 'Ba53', 'Ba51', 'Bf54', 'Qf44', 'Ka41', 'If33', 'Ea4', 'Pf51', 'Pf41', 'Pf31', 'Pf22', 'Pa21', 'Pa33', 'Pa44', 'Pa55', 'Sf42', 'Sf32', 'Sa2', 'Sa32', 'Sa43', 'Af43', 'Aa3', 'Aa42'], time:300}, log:[], type:{game:300, move:15}});
	if DEBUG {
		log.Println("LoadMatch")
	}
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
	if DEBUG {
		log.Println("ListMatches")
	}

	user, _ := user.FindUser(r)
	var matches []store.Match
	var open []store.Match

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}
	store.GetDB().Table("matches").Where("white_player_id = ? OR black_player_id = ?", user.ID, user.ID).Find(&matches)
	store.GetDB().Table("matches").Where("white_player_id != ? AND black_player_id != ? AND status='open'", user.ID, user.ID).Find(&open)

	savedMatches := make([][]string, 0)
	myOpenMatches := make([][]string, 0)
	openMatches := make([][]string, 0)
	readyMatches := make([][]string, 0)
	waitingMatches := make([][]string, 0)
	finishedMatches := make([][]string, 0)
	for m := range matches {
		if matches[m].BlackPlayerId == user.ID && matches[m].WhitePlayerId == user.ID {
			savedMatches = append(savedMatches, []string{fmt.Sprintf("%s:%d", matches[m].Title, matches[m].ID), ""})
		} else if matches[m].BlackPlayerId == user.ID || matches[m].WhitePlayerId == user.ID { // so only one player is user
			if matches[m].Game.Status == "open" {
				myOpenMatches = append(myOpenMatches, []string{fmt.Sprintf("%s:%d", matches[m].Title, matches[m].ID), ""})
			} else if matches[m].Game.Status == "engaged" {
				if matches[m].Logs == "" || len(strings.Split(strings.Trim(matches[m].Logs, " "), ":::"))%2 == 0 { // white turn
					if matches[m].WhitePlayerId == user.ID {
						readyMatches = append(readyMatches, []string{fmt.Sprintf("%s:%d", matches[m].Title, matches[m].ID), ""})
					} else {
						waitingMatches = append(waitingMatches, []string{fmt.Sprintf("%s:%d", matches[m].Title, matches[m].ID), ""})
					}
				} else { // black turn
					if matches[m].BlackPlayerId == user.ID {
						readyMatches = append(readyMatches, []string{fmt.Sprintf("%s:%d", matches[m].Title, matches[m].ID), ""})
					} else {
						waitingMatches = append(waitingMatches, []string{fmt.Sprintf("%s:%d", matches[m].Title, matches[m].ID), ""})
					}
				}
			} else if matches[m].Game.Status == "White Won" || matches[m].Game.Status == "Black Won" {
				finishedMatches = append(finishedMatches, []string{fmt.Sprintf("%s:%d", matches[m].Title, matches[m].ID), ""})
			}
		}
	}
	for m := range open {
		openMatches = append(openMatches, []string{fmt.Sprintf("%s:%d:%d", open[m].Title, open[m].ID, open[m].Game.Rating), ""})
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

func LiveMatches(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("LiveMatches")
	}

	user, _ := user.FindUser(r)

	for tk, mtch := range store.BlitzMatches() {
		if store.Online()[mtch.WhitePlayerId] == tk {
			store.Sessions()[user.Token].Inbox <- mtch
		}
	}

	utils.Respond(w, utils.Message(true, "Sending Live streams."))
}
func StartBlitz(token, game string) {
	if DEBUG {
		log.Println("StartBlitz")
	}
	gameLength := 0
	moveLength := 0
	switch game {
	case "3/900":
		gameLength = 900
		moveLength = 3
	case "5/600":
		gameLength = 600
		moveLength = 5
	case "15/360":
		gameLength = 360
		moveLength = 15
	case "10/60":
		gameLength = 60
		moveLength = 10
	case "720":
		gameLength = 720
	case "1500":
		gameLength = 1500
	}
	user := store.Sessions()[token].User
	//delete(store.BlitzMatches(), token)
	for tk, mtch := range store.BlitzMatches() {
		if DEBUG {
			log.Printf("   tk:%v\n    match:%v\n", tk, mtch)
		}
		if mtch.BlackPlayerId == 0 && mtch.Game.GameClock == uint32(gameLength) && mtch.Game.MoveClock == uint32(moveLength) {
			if DEBUG {
				log.Printf("      same game type\n")
				log.Printf("      ratings: %d vs %d\n", user.Rating(game), store.Sessions()[tk].User.Rating(game))
			}
			if store.Sessions()[tk].User.Rating(game) > user.Rating(game)-200 && store.Sessions()[tk].User.Rating(game) < user.Rating(game)+200 {
				if DEBUG {
					log.Printf("         acceptable rankings\n")
				}
				store.BlitzMap[token] = mtch
				mtch.BlackPlayerId = user.ID
				mtch.Black.UserId = mtch.BlackPlayerId
				mtch.CreatedAt = time.Now()
				mtch.UpdatedAt = time.Now()
				mtch.Title = mtch.Title + user.UserId
				store.Sessions()[tk].Inbox <- mtch
				store.Sessions()[token].Inbox <- mtch
				return
			}
		}
	}
	blitz := store.Match{Title: user.UserId + ":", Log: []string{}, WhitePlayerId: user.ID, Game: store.Type{Name: fmt.Sprintf("Blitz-%d%d", moveLength, gameLength), GameClock: uint32(gameLength), MoveClock: uint32(moveLength), Status: "Waiting"},
		White: store.Army{UserId: user.ID, Pieces: []string{"Kc44", "Qd41", "Id31", "Ed4", "Rd54", "Rd5", "Rc52", "Nd53", "Nd51", "Nc33", "Bc53", "Bc55", "Bd52", "Ad42", "Ad3", "Ac43", "Pd55", "Pd44", "Pd33", "Pd21", "Pc22", "Pc31", "Pc41", "Pc51", "Sd43", "Sd32", "Sd2", "Sc32", "Sc42"}, Time: gameLength},
		Black: store.Army{UserId: 0, Pieces: []string{"Ka41", "Qf44", "If33", "Ea4", "Ra5", "Rf52", "Ra54", "Nf53", "Nf55", "Na31", "Ba53", "Ba51", "Bf54", "Af43", "Aa3", "Aa42", "Pf51", "Pf41", "Pf31", "Pf22", "Pa21", "Pa33", "Pa44", "Pa55", "Sf42", "Sf32", "Sa2", "Sa32", "Sa43"}, Time: gameLength}}

	store.BlitzMatches()[token] = &blitz
}
func AbortBlitz(token string) {
	if DEBUG {
		log.Println("AbortBlitz")
	}
	delete(store.BlitzMap, token)
}
func BlitzMove(token, move string) {
	if DEBUG {
		log.Printf("BlitzMove (%s, %s)\n", token, move)
	}
	//user := store.Sessions()[token].User
	blitz := store.BlitzMatches()[token]
	if blitz == nil {
		store.Sessions()[token].Inbox <- "type||error|||message||blitz not found"
		return
	}
	// time management
	elapsed := time.Now().Sub(blitz.UpdatedAt)
	log.Printf("   elapsed %3f\n", elapsed.Seconds())
	if len(blitz.Log)%2 == 0 { // white
		if elapsed.Seconds() > float64(blitz.White.Time+int(blitz.Game.MoveClock)) {
			// lost on time
			blitz.WinLoseDraw("lose", "White Lost on Time")
			store.Sessions()[store.Online()[blitz.WhitePlayerId]].Inbox <- fmt.Sprintf("type||loss|||info||%s|||rating||%d", blitz.Game.Status, store.GetUser(blitz.WhitePlayerId).Rating(blitz.Game.Name))
			store.Sessions()[store.Online()[blitz.BlackPlayerId]].Inbox <- fmt.Sprintf("type||win|||info||%s|||rating||%d", blitz.Game.Status, store.GetUser(blitz.BlackPlayerId).Rating(blitz.Game.Name))
			delete(store.BlitzMatches(), store.Online()[blitz.WhitePlayerId])
			delete(store.BlitzMatches(), store.Online()[blitz.BlackPlayerId])
			return
		} else if elapsed.Seconds() > float64(blitz.Game.MoveClock) {
			// took all of move clock... remove time
			blitz.White.Time += int(blitz.Game.MoveClock) - int(elapsed.Seconds())
		} else {
			// get half of remaining time added to game clock
			blitz.White.Time += int((float64(blitz.Game.MoveClock) - elapsed.Seconds()) / 2)
		}
	} else { // black
		if elapsed.Seconds() > float64(blitz.Black.Time+int(blitz.Game.MoveClock)) {
			// lost on time
			blitz.WinLoseDraw("win", "Black Lost on Time")
			store.Sessions()[store.Online()[blitz.WhitePlayerId]].Inbox <- fmt.Sprintf("type||win|||info||%s|||rating||%d", blitz.Game.Status, store.GetUser(blitz.WhitePlayerId).Rating(blitz.Game.Name))
			store.Sessions()[store.Online()[blitz.BlackPlayerId]].Inbox <- fmt.Sprintf("type||loss|||info||%s|||rating||%d", blitz.Game.Status, store.GetUser(blitz.BlackPlayerId).Rating(blitz.Game.Name))
			delete(store.BlitzMatches(), store.Online()[blitz.WhitePlayerId])
			delete(store.BlitzMatches(), store.Online()[blitz.BlackPlayerId])
			return
		} else if elapsed.Seconds() > float64(blitz.Game.MoveClock) {
			// took all of move clock... remove time
			blitz.Black.Time += int(blitz.Game.MoveClock) - int(elapsed.Seconds())
		} else {
			// get half of remaining time added to game clock
			blitz.Black.Time += int((float64(blitz.Game.MoveClock) - elapsed.Seconds()) / 2)
		}
	}
	blitz.UpdatedAt = time.Now()
	blitz.Move(strings.Replace(move, "x", "~", 1), true)
	yourturn := blitz.BlackPlayerId
	if len(blitz.Log)%2 == 0 {
		yourturn = blitz.WhitePlayerId
	}
	yourToken := store.Online()[yourturn]
	store.Sessions()[yourToken].Inbox <- blitz
}
func BlitzEnd(token, reason string) {
	if DEBUG {
		log.Printf("BlitzEnd (%s)\n", token)
	}
	//user := store.Sessions()[token].User
	blitz := store.BlitzMatches()[token]

	if blitz == nil {
		return
	}
	blackPlayer := store.GetUser(blitz.BlackPlayerId)
	whitePlayer := store.GetUser(blitz.WhitePlayerId)
	if blackPlayer == nil || whitePlayer == nil {
		delete(store.BlitzMatches(), token)
		return
	}
	blackType := "win"
	whiteType := "loss"
	if store.Online()[blitz.WhitePlayerId] == token {
		blitz.WinLoseDraw("win", "White "+reason)
	} else if store.Online()[blitz.BlackPlayerId] == token {
		blitz.WinLoseDraw("lose", "Black "+reason)
		blackType = "loss"
		whiteType = "win"
	} else {
		// this shouldn't be the case.... hacking?
		return
	}
	store.Sessions()[store.Online()[blitz.WhitePlayerId]].Inbox <- fmt.Sprintf("type||%s|||info||%s|||rating||%d", whiteType, blitz.Game.Status, whitePlayer.Rating(blitz.Game.Name))
	store.Sessions()[store.Online()[blitz.BlackPlayerId]].Inbox <- fmt.Sprintf("type||%s|||info||%s|||rating||%d", blackType, blitz.Game.Status, blackPlayer.Rating(blitz.Game.Name))
	delete(store.BlitzMatches(), store.Online()[blitz.WhitePlayerId])
	delete(store.BlitzMatches(), store.Online()[blitz.BlackPlayerId])
}
