package hold

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	u "Chexx/server/utils"

	"github.com/jinzhu/gorm"
)

type GameState int

const (
	CHALLANGE = iota
	WHITE2MOVE
	BLACK2MOVE
	WHITEWONONTIME
	BLACKWONONTIME
	WHITEABDICATED
	BLACKABDICATED
	WHITEWONBYCHECKMATE
	BLACKWONBYCHECKMATE
	STALEMATE
	PRACTICE
	PUZZLE
	DELETED
	LASTSTATE
) // Game.ts must be synced

type GameType int

const (
	BLITZ5SECOND = iota
	BLITZ20SECOND
	GAME10MINUTE
	GAME30MINUTE
	FORMAL2DAY
	FORMAL3DAY
	FORMAL5DAY
	TOURNAMENT15MIN
	TOURNAMENT1HOUR
	TOURNAMENT4HOUR
	LASTTYPE
)

type GameFormation int

const (
	RAID       = iota // = 4
	CLASH             // = 5
	STANDARD          // = 6
	MIGHTY            // = 7
	GRAND             // = 8
	MASSIVE           // = 9
	COLOSSAL          // = 10
	TITANIC           // = 12
	MONUMENTAL        // = 14
	EPIC              // = 17
	THEATRE           // = 21
	LASTFORMATION
)

type Game struct {
	White     string `json:"white"`
	Black     string `json:"black"`
	WhiteTurn bool   `json:"whiteTurn"`
	Board     int    `json:"board"`
	Name      string `json:"name"`
	Quest     string `json:"quest"`
	Depth     int    `json:"depth"`
	Move      string
	Score     float64
	AI        string
	Zwanzig   []string
	ZwanPtr   int
}

func (game *Game) String() string {
	st := strings.TrimSpace(game.White) + "  ///  " + strings.TrimSpace(game.Black) + "\n"
	st += fmt.Sprintf("Board: %d    White Move? %v   Name: %s   Quest: %s   AI: %s\n", game.Board, game.WhiteTurn, game.Name, game.Quest, game.AI)
	st += fmt.Sprintf("Move: %s", game.Move)
	return st
}

func (game *Game) ZwanzigTest(move string) bool {
	if game.Zwanzig == nil || game.Zwanzig[0] == "" || game.Zwanzig[1] == "" || game.Zwanzig[2] == "" || game.Zwanzig[3] == "" ||
		game.Zwanzig[4] == "" || game.Zwanzig[5] == "" || game.Zwanzig[6] == "" || game.Zwanzig[7] == "" {
		return false
	}
	m1 := strings.Split(game.Zwanzig[game.ZwanPtr], "~")
	m2 := strings.Split(game.Zwanzig[(game.ZwanPtr+1)%4], "~")
	m3 := strings.Split(game.Zwanzig[(game.ZwanPtr+2)%4], "~")
	m4 := strings.Split(game.Zwanzig[(game.ZwanPtr+3)%4], "~")
	m5 := strings.Split(game.Zwanzig[(game.ZwanPtr+4)%4], "~")
	m6 := strings.Split(game.Zwanzig[(game.ZwanPtr+5)%4], "~")
	m7 := strings.Split(game.Zwanzig[(game.ZwanPtr+6)%4], "~")
	m8 := strings.Split(game.Zwanzig[(game.ZwanPtr+7)%4], "~")
	mv := strings.Split(move, "~")
	if m1[0] == m7[1] && m1[1] == m7[0] && m6[0] == m8[1] && m6[1] == m8[0] && mv[0] == m1[0] && mv[1] == m1[1] {
		return true
	}
	if m1[0] == m5[1] && m5[0] == m7[1] && m7[0] == m1[1] && m8[0] == m4[1] && m4[0] == m6[1] && m6[0] == m8[1] && mv[0] == m1[0] && mv[1] == m1[1] {
		return true
	}
	if m1[0] == m3[1] && m3[0] == m5[1] && m5[0] == m7[1] && m7[0] == m1[1] &&
		m8[0] == m2[1] && m2[0] == m4[1] && m4[0] == m6[1] && m6[0] == m8[1] && mv[0] == m1[0] && mv[1] == m1[1] {
		return true
	}

	return false
}

func (game *Game) SetMove(move string) {
	if game.Zwanzig == nil {
		game.Zwanzig = make([]string, 8)
	}
	game.Zwanzig[game.ZwanPtr] = move
	game.ZwanPtr = (game.ZwanPtr + 1) % 8
}

type Player struct {
	Id        uint          `json:"id"`
	Name      string        `json:"name"`
	Pieces    string        `json:"pieces"`
	GameClock time.Duration `json:"gameClock"`
	MoveClock time.Duration `json:"moveClock"`
	Rating    int           `json:"rating"`
	Pict      string        `json:"pict"`
}

type Match struct {
	gorm.Model
	Name        string        `json:"name"`
	White       Player        `json:"white" gorm:"-"`
	WhiteId     uint          `json:"-" gorm:"whiteId"`
	WhitePieces string        `json:"-" gorm:"whitePieces"`
	WGClock     time.Duration `json:"-" gorm:"wgClock"`
	WMClock     time.Duration `json:"-" gorm:"wmClock"`
	Black       Player        `json:"black" gorm:"-"`
	BlackId     uint          `json:"-" gorm:"blackId"`
	BlackPieces string        `json:"-" gorm:"blackPieces"`
	BGClock     time.Duration `json:"-" gorm:"bgClock"`
	BMClock     time.Duration `json:"-" gorm:"bmClock"`
	History     string        `json:"history"`
	GType       GameType      `json:"gameType"`
	State       GameState     `json:"gameState"`
	Formation   GameFormation `json:"gameForm"`
	GameTime    time.Time     `json:"-"`
	Min         int           `json:"min"`
	Max         int           `json:"max"`
}

// var BlitzMatches map[WhiteId]*Match

var BlitzLoungeMux sync.Mutex
var BlitzLounge map[uint]*Match

var OnlookersMux sync.Mutex
var Onlookers map[*Match][]*Session

func (m Match) String() string {
	str := fmt.Sprintf("\n---Match(%d) %s   type:%d  formation:%d---", m.ID, m.Name, m.GType, m.Formation)
	str += "\n   -- [" + m.WhitePieces + " || " + m.BlackPieces + "]"
	str += "\n   --History: {" + m.History + "}"
	str += fmt.Sprintf("\n   --Clock: %v", m.GameTime)
	str += fmt.Sprintf("\n   --White: %v", m.White)
	str += fmt.Sprintf("\n   --Black: %v", m.Black)
	str += fmt.Sprint("\n-----------------------------------------------")
	return str
}

func (p Player) String() string {
	str := fmt.Sprintf("%s(%d):%d", p.Name, p.Id, p.Rating)
	str += fmt.Sprintf("\n    Clock: %v||%v    [%v]", p.GameClock, p.MoveClock, p.Pieces)
	return str
}

func (match *Match) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(b), &m); err != nil {
		return err
	}

	if m["name"] == nil || m["gameType"] == nil || m["gameForm"] == nil || m["min"] == nil || m["max"] == nil || m["white"] == nil || m["black"] == nil {
		return errors.New("Insufficient info for new match")
	}
	match.Name = m["name"].(string)
	match.GType = GameType(uint(m["gameType"].(float64)))
	match.Formation = GameFormation(uint(m["gameForm"].(float64)))
	match.Min = int(m["min"].(float64))
	match.Max = int(m["max"].(float64))
	wh := m["white"].(map[string]interface{})
	bl := m["black"].(map[string]interface{})
	match.White = Player{Id: uint(wh["id"].(float64)), Name: wh["name"].(string), Pieces: wh["pieces"].(string)}
	match.Black = Player{Id: uint(bl["id"].(float64)), Name: bl["name"].(string), Pieces: bl["pieces"].(string)}
	match.WhiteId = match.White.Id
	match.WhitePieces = match.White.Pieces
	match.WGClock = match.White.GameClock
	match.WMClock = match.White.MoveClock
	match.BlackId = match.Black.Id
	match.BlackPieces = match.Black.Pieces
	match.BGClock = match.Black.GameClock
	match.BMClock = match.Black.MoveClock

	return nil
}

func (p *Player) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(b), &m); err != nil {
		return err
	}
	p.Id = m["id"].(uint)
	p.Name = m["name"].(string)
	p.Pieces = m["pieces"].(string)
	p.GameClock = m["gameClock"].(time.Duration)
	p.MoveClock = m["moveClock"].(time.Duration)
	p.Rating = m["rating"].(int)

	return nil
}

func (match *Match) Challange() map[string]interface{} {

	if resp, ok := match.Validate(); !ok {
		return resp
	}

	match = match.Create()

	GetDB().Create(match)

	if match.ID <= 0 {
		return u.Message(false, "Failed to create match, connection error.")
	}

	response := u.Message(true, "Match has been created")
	response["match"] = match
	return response
}

func (match *Match) Create() *Match {
	debug := 6

	log(debug, "Create Match: "+match.String())

	switch match.Formation {
	case MASSIVE:
		break
	case COLOSSAL:
		match.WhitePieces = "Ne5 Pd44 Nd32 Pd2 Pc41 Nc5 Pd55 Pd43 Pd31 Pc33 Nc32 Pc42 Pc51 Pd66 Pd42 Pc43 Pc61 Pd77 Sd65 Sd53 Sd41 Sc44 Sc53 Sc62 Pc71 Sd76 Ad64 Ad52 Ad4 Ac54 Ac63 Sc72 Rd75 Rc73 Nd74 Bc74 Bd73 Nc75 Ed72 Bd6 Ic76 Kd71 Qc77 Rd7 Ld54 Lc52 cd3"
		match.BlackPieces = "Ra7 Kf77 Qa71 Ef76 Ia72 Bf75 Na73 Nf74 Ba6 Ba74 Rf73 Ra75 Sf72 Af63 Af54 Aa4 Aa52 Aa64 Sa76 Pf71 Sf62 Sf53 Sf44 Sa41 Sa53 Sa65 Pa77 Pf61 Pf43 Pa42 Pa66 Pf51 Pf42 Pf33 Pa31 Pa43 Na32 Pa55 Nf5 Pf41 Nf32 Pa2 Pa44 Nb5 Lf52 La54 ca3"
		break
	case GRAND:
		match.WhitePieces = "Pe6 Sd66 Rd65 Nd64 Bd63 Ed62 Kd61 Rd6 Qc66 Ic65 Nc64 Bc63 Rc62 Sc61 Pc6 Pe5 Sd55 Ad54 Bd5 Ac52 Sc51 Pc5 Pe4 Pd44 Sd43 Ad42 Ac43 Sc42 Pc41 Pc4 Pd33 Pd32 Sd31 Ad3 Sc33 Pc32 Pc31 Pd21 Pd2 Pc22"
		match.BlackPieces = "Pf6 Sf61 Rf62 Nf63 Bf64 Ef65 Kf66 Ra6 Qa61 Ia62 Na63 Ba64 Ra65 Sa66 Pb6 Pf5 Sf51 Af52 Ba5 Aa54 Sa55 Pb5 Pf4 Pf41 Sf42 Af43 Aa42 Sa43 Pa44 Pb4 Pf31 Pf32 Sf33 Aa3 Sa31 Pa32 Pa33 Pf22 Pa2 Pa21"
		break
	case MIGHTY:
		match.WhitePieces = "Pe5 Rd55 Nd54 Bd53 Ed52 Kd51 Rd5 Qc55 Ic54 Nc53 Bc52 Rc51 Pc5 Pe4 Sd44 Ad43 Bd4 Ac42 Sc41 Pc4 Pe3 Pd33 Sd32 Ad31 Ac33 Sc32 Pc31 Pc3 Pd22 Pd21 Sd2 Pc22 Pc21 Nd1"
		match.BlackPieces = "Pf5 Rf51 Nf52 Bf53 Ef54 Kf55 Ra5 Qa51 Ia52 Na53 Ba54 Ra55 Pb5 Pf4 Sf41 Af42 Ba4 Aa43 Sa44 Pb4 Pf3 Pf31 Sf32 Af33 Aa31 Sa32 Pa33 Pb3 Pf21 Pf22 Sa2 Pa21 P22 Na1"
		break
	case STANDARD:
		match.WhitePieces = "Kd4 Qc44 Bc43 Nc42 Rc41 Pc4 Bd41 Nd42 Bd43 Rd44 Pe4 Pc3 Pc31 Sc32 Ac33 Ad3 Ad31 Sd32 Pd33 Pe3 Pc21 Pc22 Sd2 Pd21 Pd22 Nd1"
		match.BlackPieces = "Ka4 Qa41 Ba42 Na43 Ra44 Pb4 Pf4 Rf41 Bf42 Nf43 Bf44 Pf3 Pf31 Sf32 Af33 Aa3 Aa31 Sa32 Pa33 Pb3 Pf21 Pf22 Sa2 Pa21 Pa22 Na1"
		break
	case CLASH:
		match.WhitePieces = "Kd3 Bc33 Nc32 Rc31 Pc3 Nd31 Bd32 Rd33 Pe3 Pc2 Pc21 Nc22 Qd2 Bd21 Pd22 Pe2 Pc11 Pd1 Pd11"
		match.BlackPieces = "Pb3 Ra33 Na32 Ba31 Ka3 Nf33 Bf32 Rf31 Pf3 Pb2 Pa22 Na21 Qa2 Bf22 Pf21 Pf2 Pa11 Pa1 Pf11"
		break
	case RAID:
		match.WhitePieces = "Kd2 Bc22 Rc21 Pc2 Pe2 Rd22 Nd21 Pe1 Bd11 Qd1 Nc11 Pc1 Pe Pd Pc"
		match.BlackPieces = "Pf2 Rf21 Nf22 Ka2 Ba21 Ra22 Pb2 Pf1 Bf11 Qa1 Na11 Pb1 Pf Pa Pb"
		break
	}

	switch match.GType {
	case BLITZ5SECOND:
		match.WGClock = time.Duration(0)
		match.WMClock, _ = time.ParseDuration("5s")
		break
	case BLITZ20SECOND:
		match.WGClock = time.Duration(0)
		match.WMClock, _ = time.ParseDuration("20s")
		break
	case GAME10MINUTE:
		match.WGClock, _ = time.ParseDuration("10m")
		match.WMClock = time.Duration(0)
		break
	case GAME30MINUTE:
		match.WGClock, _ = time.ParseDuration("30m")
		match.WMClock = time.Duration(0)
		break
	case FORMAL2DAY:
		match.WGClock, _ = time.ParseDuration("72h")
		match.WMClock, _ = time.ParseDuration("48h")
		break
	case FORMAL3DAY:
		match.WGClock, _ = time.ParseDuration("96h")
		match.WMClock, _ = time.ParseDuration("72h")
		break
	case FORMAL5DAY:
		match.WGClock, _ = time.ParseDuration("168h")
		match.WMClock, _ = time.ParseDuration("120h")
		break
	case TOURNAMENT15MIN:
		match.WGClock, _ = time.ParseDuration("15m")
		match.WMClock, _ = time.ParseDuration("10s")
		break
	case TOURNAMENT1HOUR:
		match.WGClock, _ = time.ParseDuration("1h")
		match.WMClock, _ = time.ParseDuration("10s")
		break
	case TOURNAMENT4HOUR:
		match.WGClock, _ = time.ParseDuration("4h")
		match.WMClock, _ = time.ParseDuration("10s")
		break
	}
	match.BGClock = match.WGClock
	match.BMClock = match.WMClock

	match.GameTime = time.Now()

	return match
}

func (match *Match) Validate() (map[string]interface{}, bool) {

	match.State = CHALLANGE
	if match.White.Id == 0 && match.Black.Id == 0 {
		return u.Message(false, "No challenger"), false
	}
	if len(match.Name) < 5 {
		return u.Message(false, "Match name is too short"), false
	}

	//Name must be unique
	temp := &Match{}

	//check for errors and duplicate emails
	err := GetDB().Table("matches").Where("name = ?", match.Name).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Name != "" {
		return u.Message(false, "Name is already in use by another match"), false
	}

	return u.Message(false, "Requirement passed"), true
}

func GetMatch(u uint, pop bool) *Match {
	match := &Match{}
	GetDB().Table("matches").Where("id = ?", u).First(match)
	if match.Name == "" { //Match not found!
		return nil
	}
	if pop {
		match.Users()
	}

	return match
}

func (match *Match) Users() {
	if match.WhiteId > 0 {
		white := GetUser(match.WhiteId)
		match.White.Name = white.Name
		match.White.Id = match.WhiteId
		match.White.Pict = white.Pict
		match.White.Rating = white.GetRating(match)
	} else {
		match.White.Name = "Pending"
	}
	if match.BlackId > 0 {
		black := GetUser(match.BlackId)
		match.Black.Name = black.Name
		match.Black.Id = match.BlackId
		match.Black.Pict = black.Pict
		match.Black.Rating = black.GetRating(match)
	} else {
		match.Black.Name = "Pending"
	}
	match.White.Pieces = match.WhitePieces
	elapsed := time.Now().Sub(match.GameTime)
	match.White.GameClock = match.WGClock
	if match.State == WHITE2MOVE {
		match.White.MoveClock = match.WMClock - elapsed
	} else {
		match.White.MoveClock = match.WMClock
	}
	match.Black.Pieces = match.BlackPieces
	match.Black.GameClock = match.BGClock
	if match.State == BLACK2MOVE {
		match.Black.MoveClock = match.BMClock - elapsed
	} else {
		match.Black.MoveClock = match.BMClock
	}
	match.GameTime = time.Now()
}

func GetPuzzle(n string) *Match {
	match := &Match{}
	GetDB().Table("matches").Where("name = ?", n).First(match)
	if match.Name == "" { //Match not found!
		return nil
	}
	match.Users()
	return match
}

func GetPuzzles(low, high int) []string {
	var names []string
	db.Model(&Match{}).Pluck("name", &names)
	db.Table("matches").Where("min > low and min < high", low, high).Pluck("name", &names)
	return names
}

func GetSaved(id uint) []string {
	var names []string
	db.Model(&Match{}).Pluck("name", &names)
	db.Table("matches").Where("black_id = ? AND white_id = ?", id, id).Pluck("name", &names)
	return names
}

func (match *Match) Update() map[string]interface{} {

	match.GameTime = time.Now()

	GetDB().Save(match)

	response := u.Message(true, "Match has been updated")
	response["match"] = match
	return response
}

func GetOpenChallanges(userId uint) []Match {

	var matches, filtered []Match
	user := &Account{}
	GetDB().Table("accounts").Select("ranks").Where("ID = ?", userId).First(&user)
	ranks := strings.Split(user.Ranks, "|")

	GetDB().Table("matches").Where("black_id != ? AND white_id != ? AND (black_id = 0 OR white_id = 0)", userId, userId).Find(&matches)
	if len(matches) == 0 { //No matches found
		return nil
	}
	for _, m := range matches {
		userRank, err := strconv.Atoi(ranks[int(m.Formation)*10+int(m.GType)])
		if err == nil && userRank >= m.Min && userRank <= m.Max {
			filtered = append(filtered, m)
		}
	}
	populate(filtered)
	return filtered
}

func GetActiveMatches(userId uint) []Match {

	var matches []Match

	GetDB().Table("matches").Where("(black_id = ? OR white_id = ?) AND black_id != 0 AND white_id != 0 AND state in (0,1,2)", userId, userId).Find(&matches)
	if len(matches) == 0 { //No matches found
		return nil
	}

	populate(matches)
	return matches
}

func GetMyChallanges(userId uint) []Match {

	var matches []Match

	GetDB().Table("matches").Where("(black_id = ? OR white_id = ?) AND (black_id = 0 OR white_id = 0) AND state = 0", userId, userId).Find(&matches)
	if len(matches) == 0 { //No matches found
		return nil
	}

	populate(matches)
	return matches
}

func AcceptChallange(userId, matchId uint) *Match {

	user := &Account{}
	GetDB().Table("accounts").Where("ID = ?", userId).First(&user)

	match := GetMatch(matchId, true)
	if match.WhiteId == 0 {
		match.WhiteId = userId
		match.White.Id = userId
		match.White.Name = user.Name
	}
	if match.BlackId == 0 {
		match.BlackId = userId
		match.Black.Id = userId
		match.Black.Name = user.Name
	}

	if err := GetDB().Save(&match).Error; err != nil {
		fmt.Println(err)
		return nil
	}
	return match
}

func populate(matches []Match) {
	for idx, _ := range matches {
		matches[idx].Users()
	}
}

func (match *Match) Timeout() bool {
	debug := 6
	log(debug, fmt.Sprintf("<match::Timeout>\n%v", match))
	elapsed := time.Now().Sub(match.GameTime)
	if match.State == CHALLANGE {
		log(debug, fmt.Sprintf("</match::Timeout> CHALLANGE"))
		match.GameTime = time.Now()
		return false
	}
	if match.State == WHITE2MOVE {
		log(debug, fmt.Sprintf("WHITE2MOVE"))
		if match.White.MoveClock < elapsed {
			log(debug, fmt.Sprintf("white clock expired"))
			match.White.GameClock = match.White.GameClock - elapsed + match.White.MoveClock
			if match.White.GameClock < 0 {
				log(debug, fmt.Sprintf("white game expired"))
				match.State = BLACKWONONTIME
				log(debug, fmt.Sprintf("</match::Timeout> BLACKWONONTIME"))
				return true
			}
		}
	}
	if match.State == BLACK2MOVE {
		log(debug, fmt.Sprintf("BLACK2MOVE"))
		if match.Black.MoveClock < elapsed {
			log(debug, fmt.Sprintf("black clock expired"))
			match.Black.GameClock = match.Black.GameClock - elapsed + match.Black.MoveClock
			if match.Black.GameClock < 0 {
				log(debug, fmt.Sprintf("black game expired"))
				match.State = WHITEWONONTIME
				log(debug, fmt.Sprintf("</match::Timeout> WHITEWONONTIME"))
				return true
			}
		}
	}
	match.GameTime = time.Now()
	log(debug, fmt.Sprintf("</match::Timeout>"))
	return false
}

func (match *Match) MoveTime() time.Duration {
	var timer time.Duration
	switch match.GType {
	case BLITZ5SECOND:
		timer, _ = time.ParseDuration("5s")
	case BLITZ20SECOND:
		timer, _ = time.ParseDuration("20s")
	case GAME10MINUTE:
	case GAME30MINUTE:
	case FORMAL2DAY:
		timer, _ = time.ParseDuration("48h")
	case FORMAL3DAY:
		timer, _ = time.ParseDuration("72h")
	case FORMAL5DAY:
		timer, _ = time.ParseDuration("96h")
	case TOURNAMENT15MIN:
	case TOURNAMENT1HOUR:
	case TOURNAMENT4HOUR:
	}
	return timer
}
