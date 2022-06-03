package store

import (
	"github.com/jarrancarr/ChexxServer/utils"
	"github.com/jinzhu/gorm"
)

type Army struct {
	gorm.Model
	Player string   `json:"id"`
	Pieces []string `json:"pieces"`
	Time   int32    `json:"time"`
}

type Type struct {
	gorm.Model
	Name      string `json:"name"`
	GameClock int32  `json:"game"`
	MoveClock int32  `json:"move"`
}
type Match struct {
	gorm.Model
	Title  string   `json:"name"`
	White  Army     `json:"white"`
	UserId string   `json:"userID"`
	Log    []string `json:"log" gorm:"-"`
	Logs   string   `json:"-"`
	Blog   string   `json:"blog"`
	Game   Type
}

type MatchWrap struct {
	Match Match `json:"match"`
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

	return m
}
