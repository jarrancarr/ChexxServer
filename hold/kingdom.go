package hold

import (
	"fmt"

	u "Chexx/server/utils"

	"github.com/jinzhu/gorm"
)

type Kingdom struct {
	gorm.Model
	Name     string `json:"name"`
	UserId   uint
	LordId   uint              `json:"lord"`
	Flag     string            `json:"flag"`
	World    string            `json:"world"`
	Lon      int32             `json:"lon"` // 40,000km circumference
	Lat      int32             `json:"lat"`
	Topology []uint32          `json:"topology" gorm:"-"`
	Roads    [][]uint32        `json:"roads" gorm:"-"`
	Inv      map[string]string `json:"inv" gorm:"-"`
}

func (k *Kingdom) String() string {
	return fmt.Sprintf("Kingdom: %s of %s\n   lat:%d    log:%d", k.Name, k.World, k.Lat, k.Lon)
}

func (k *Kingdom) Validate() (map[string]interface{}, bool) {

	if k.UserId < 1 {
		return u.Message(false, "No User"), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (k *Kingdom) Create() map[string]interface{} {

	if resp, ok := k.Validate(); !ok {
		return resp
	}

	GetDB().Create(k)

	if k.ID <= 0 {
		return u.Message(false, "Failed to create kingdom, connection error.")
	}

	response := u.Message(true, "Kingdom created")
	response["id"] = k.ID
	return response
}

func (k *Kingdom) Update() map[string]interface{} {
	GetDB().Save(k)
	response := u.Message(true, "Kingdom has been updated")
	return response
}

func GetUserKingdom(id uint) *Kingdom {
	debug := 1
	log(debug, fmt.Sprintf("<kingdom::GetUserKingdom>%v\n", id))
	kingdom := &Kingdom{}
	GetDB().Table("kingdoms").Where("user_id = ?", id).First(kingdom)
	log(debug, fmt.Sprintf("</kingdom> %v\n", kingdom))
	if kingdom.Name == "" {
		return nil
	}
	return kingdom
}
