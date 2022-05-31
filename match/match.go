package match

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	Game   Type
}

type MatchWrap struct {
	Match Match `json:"match"`
}

func Matches(w http.ResponseWriter, r *http.Request) {

	var m MatchWrap
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
