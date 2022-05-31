package hold

import (
	"fmt"
	"strconv"
	"strings"

	u "Chexx/server/utils"

	"github.com/jinzhu/gorm"
)

type AI struct {
	gorm.Model
	Name         string             `gorm:"name"`
	Piece        map[string]float64 `gorm:"-"`
	PieceValues  string             `gorm:"pieces"`
	Attack       map[string]float64 `gorm:"-"`
	AttackValues string             `gorm:"attacks"`
	Defend       map[string]float64 `gorm:"-"`
	DefendValues string             `gorm:"defend"`
	Space        map[string]float64 `gorm:"-"`
	SpaceValues  string             `gorm:"spaces"`
	Depth        int                `gorm:"depth"`
	Consider     int                `gorm:"consider"` // how many options to consider another level
	Diminish     int                `gorm:"diminish"`
	Wins         int                `gorm:"-"`
	Losses       int                `gorm:"-"`
	Draws        int                `gorm:"-"`
	Points       int                `gorm:"points"`
}

func (ai *AI) String() string {
	st := fmt.Sprintf("Name: %s   Depth: %d  Consider: %d  Diminish: %d  Points: %d\n", ai.Name, ai.Depth, ai.Consider, ai.Diminish, ai.Points)
	st += fmt.Sprintf("  Material:   P/p/S/s: %2.3f/%2.3f/%2.3f/%2.3f\n", ai.Piece["P"], ai.Piece["p"], ai.Piece["S"], ai.Piece["s"])
	st += fmt.Sprintf("              N/B/A/R/I/E/Q/K: %2.3f/%2.3f/%2.3f/%2.3f/%2.3f/%2.3f/%2.3f/%2.3f\n",
		ai.Piece["N"], ai.Piece["B"], ai.Piece["A"], ai.Piece["R"], ai.Piece["I"], ai.Piece["E"], ai.Piece["Q"], ai.Piece["K"])
	return st
}

func (ai *AI) Validate() (map[string]interface{}, bool) {
	if ai.Name == "" {
		return u.Message(false, "No Name"), false
	}
	if ai.Depth*ai.Diminish >= ai.Consider {
		return u.Message(false, "Not enough to consider"), false
	}
	if len(ai.Piece) < 16 || len(ai.Space) < 16 || len(ai.Attack) < 16*16 || len(ai.Defend) < 16*16 {
		return u.Message(false, "Undefined pieces"), false
	}
	return u.Message(false, "Requirement passed"), true
}

func (ai *AI) Create() (map[string]interface{}, bool) {
	if resp, ok := ai.Validate(); !ok {
		return resp, false
	}
	// convert maps to strings
	ai.AttackValues = convert(ai.Attack)
	ai.PieceValues = convert(ai.Piece)
	ai.DefendValues = convert(ai.Defend)
	ai.SpaceValues = convert(ai.Space)
	GetDB().Create(ai)
	return u.Message(false, "Brain created"), true
}

func (ai *AI) Update() map[string]interface{} {
	GetDB().Save(ai)
	response := u.Message(true, "AI has been updated")
	return response
}

func convert(data map[string]float64) string {
	conv := ""
	for d := range data {
		if conv != "" {
			conv += "|"
		}
		conv += d + ":" + fmt.Sprintf("%f", data[d])
	}
	return conv
}

func revert(data string) map[string]float64 {
	conv := make(map[string]float64)
	for _, d := range strings.Split(data, "|") {
		c := strings.Split(d, ":")
		conv[c[0]], _ = strconv.ParseFloat(c[1], 64)
	}
	return conv
}

func GetAI(name string) *AI {
	ai := &AI{}
	GetDB().Table("ais").Where("name = ?", name).First(ai)
	if ai.Name == "" { //AI not found!
		return nil
	}
	ai.Attack = revert(ai.AttackValues)
	ai.Piece = revert(ai.PieceValues)
	ai.Defend = revert(ai.DefendValues)
	ai.Space = revert(ai.SpaceValues)
	return ai
}

func GetAIsAvailable() []string {
	var names []string
	db.Model(&AI{}).Pluck("name", &names)
	db.Table("ais").Where("1=1").Pluck("name", &names)
	return names
}

func GetAllAIs() []*AI {
	var ais []*AI
	GetDB().Table("ais").Where("deleted_at is not null").Find(&ais)
	return ais
}

func DeleteAI(id uint) {
	GetDB().Table("ais").Delete(AI{}, "ID =?", id)
}
