package store

import (
	"encoding/json"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/jarrancarr/ChexxServer/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type About struct {
	gorm.Model
	School    string `json:"school"`
	City      string `json:"city"`
	Company   string `json:"company"`
	DOB       uint   `json:"dob"`
	Origin    string `json:"origin"`
	Profile   string `json:"profile"`
	Interests string `json:"interests"`
}

type Friend struct {
	UserId string
}
type User struct {
	gorm.Model
	Token         string            `json:"token" gorm:"-"`
	Expire        int64             `json:"data_access_expiration_time"`
	Email         string            `json:"email"`
	UserId        string            `json:"userid"`
	Name          string            `json:"fullName"`
	Password      string            `json:"password"`
	Prop          map[string]string `json:"property" gorm:"-"`
	Property      string            `json:"-"`
	SignedRequest string            `json:"signedRequest"`
	Rank          uint32            `json:"rank"`
	About         About             `json:"about"`
	Friend        []string          `json:"friend" gorm:"-"`
	Friends       string            `json:"-" gorm:"friends"`
	Hangout       []string          `json:"hangout" gorm:"-"`
	Hangouts      string            `json:"-" gorm:"hangouts"`
	// Picture       []uint8 `json:"picture"`
}

type Session struct {
	User *User
	// NotificationQueue []*Comment
	NumNewMoves int
	Blitz       *Match
	//Watching             *Match
	//Polling           bool
	WsConn *websocket.Conn
	Inbox  chan interface{}
}

var SessionMap map[string]*Session // map of tokens to sessions
var OnlineMapping map[uint]string  // map of ids to tokens

func Sessions() map[string]*Session {
	if SessionMap == nil {
		SessionMap = make(map[string]*Session)
	}
	return SessionMap
}
func Online() map[uint]string {
	if OnlineMapping == nil {
		OnlineMapping = make(map[uint]string)
	}
	return OnlineMapping
}

type Team struct {
	gorm.Model
	Name      string           `json:"name"`
	Officers  map[string]*User `gorm:"-"`
	Members   map[string]*User `gorm:"-"`
	Applicant map[string]*User `gorm:"-"`
	About     string           `json:"about"`
	Pict      string           `json:"pict"`
}

func (u *User) Create() map[string]interface{} {

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)
	u.Rank = 1000

	GetDB().Create(u)

	if u.ID <= 0 {
		return utils.Message(false, "Failed to create account, connection error.")
	}

	tk := &Token{UserId: u.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	u.Token = tokenString

	u.Password = "" //delete password

	response := utils.Message(true, "Account has been created")
	response["user"] = u
	return response
}
func GetUser(id uint) *User {

	u := &User{}
	GetDB().Table("users").Where("id = ?", id).First(u)
	if u.Name == "" || u.Email == "" { //User not found!
		return nil
	}
	json.Unmarshal([]byte(u.Property), &u.Prop)
	if u.Prop == nil {
		u.Prop = make(map[string]string)
		u.Prop["test"] = "success"
	}
	json.Unmarshal([]byte(u.Friends), &u.Friend)
	if u.Friend == nil {
		u.Friend = []string{}
	}
	json.Unmarshal([]byte(u.Hangouts), &u.Hangout)
	if u.Hangout == nil {
		u.Hangout = []string{"International Lounge"}
	}
	return u
}
func (u *User) Validate() (map[string]interface{}, bool) {

	if u.Email == "" {
		return utils.Message(false, "No Email"), false
	}
	if u.Name == "" {
		return utils.Message(false, "No Name"), false
	}
	if u.UserId == "" {
		return utils.Message(false, "No UserId"), false
	}

	return utils.Message(true, "Requirement passed"), true
}
func (u *User) Update() map[string]interface{} {

	if resp, ok := u.Validate(); !ok {
		return resp
	}

	// convert properties

	prop, err := json.Marshal(u.Prop)
	if err == nil {
		u.Property = string(prop)
	}
	friends, err := json.Marshal(u.Friend)
	if err == nil {
		u.Friends = string(friends)
	}
	hangouts, err := json.Marshal(u.Hangout)
	if err == nil {
		u.Hangouts = string(hangouts)
	}

	GetDB().Save(u)

	if u.ID <= 0 {
		return utils.Message(false, "Failed to create user, connection error.")
	}

	response := utils.Message(true, "User updated")
	return response
}
