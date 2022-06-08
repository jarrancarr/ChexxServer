package store

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jarrancarr/ChexxServer/utils"
	"github.com/jinzhu/gorm"
)

type Army struct {
	UserId uint     `json:"userid"`
	Pieces []string `json:"pieces"`
	Time   int      `json:"time"`
}

type Type struct {
	Name      string `json:"name"`
	GameClock uint32 `json:"game"`
	MoveClock uint32 `json:"move"`
}
type Match struct {
	gorm.Model
	Title         string   `json:"name"`
	White         Army     `json:"white" gorm:"-"`
	Black         Army     `json:"black" gorm:"-"`
	WhiteArmy     string   `json:"-"`
	BlackArmy     string   `json:"-"`
	WhitePlayerId uint     `json:"-"`
	BlackPlayerId uint     `json:"-"`
	Log           []string `json:"log" gorm:"-"`
	Logs          string   `json:"-"`
	Blog          string   `json:"blog"`
	Game          Type     `gorm:"embedded"`
}

func (m *Match) Validate() (map[string]interface{}, bool) {

	if m.Title == "" {
		return utils.Message(false, "No Title"), false
	}

	return utils.Message(false, "Requirement passed"), true
}

func (m *Match) Create() map[string]interface{} {

	if resp, ok := m.Validate(); !ok {
		return resp
	}

	m.Logs = strings.Join(m.Log, " ")
	m.BlackArmy = strings.Join(m.Black.Pieces, " ") + "|" + fmt.Sprintf("%d", m.Black.Time)
	m.WhiteArmy = strings.Join(m.White.Pieces, " ") + "|" + fmt.Sprintf("%d", m.White.Time)
	m.BlackPlayerId = m.Black.UserId
	m.WhitePlayerId = m.White.UserId

	GetDB().Create(m)

	if m.ID <= 0 {
		return utils.Message(false, "Failed to create message, connection error.")
	}

	response := utils.Message(true, "Match created")
	response["id"] = m.ID
	return response
}

func GetMatch(id uint) *Match {

	m := &Match{}
	GetDB().Table("matches").Where("id = ?", id).First(m)

	m.Log = strings.Split(m.Logs, " ")

	black := strings.Split(m.BlackArmy, "|")
	blackClock, _ := strconv.Atoi(black[1])
	m.Black = Army{UserId: m.BlackPlayerId, Pieces: strings.Split(black[0], " "), Time: blackClock}

	white := strings.Split(m.WhiteArmy, "|")
	whiteClock, _ := strconv.Atoi(white[1])
	m.White = Army{UserId: m.WhitePlayerId, Pieces: strings.Split(white[0], " "), Time: whiteClock}

	return m
}
