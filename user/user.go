package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jarrancarr/ChexxServer/utils"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Token  string `json:"accessToken"`
	Expire int64  `json:"data_access_expiration_time"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	// Picture       []uint8 `json:"picture"`
	SignedRequest string `json:"signedRequest"`
	UserId        string `json:"userID"`
}

type UserWrap struct {
	Guy User `json:"thisGuy"`
}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	var u UserWrap

	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}
	fmt.Println("------------------------------")
	fmt.Println(u.Guy.Name)
	fmt.Println("------------------------------")
	utils.CorsHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("good"))
}
