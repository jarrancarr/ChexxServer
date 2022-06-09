package store

import (
	"os"

	"github.com/golang-jwt/jwt"
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

type User struct {
	gorm.Model
	Token         string            `json:"token" gorm:"-"`
	Expire        int64             `json:"data_access_expiration_time"`
	Email         string            `json:"email"`
	Name          string            `json:"fullName"`
	Password      string            `json:"password"`
	Property      map[string]string `json:"property" gorm:"-"`
	SignedRequest string            `json:"signedRequest"`
	UserId        string            `json:"userid"`
	Rank          uint32            `json:"rank"`
	About         About             `json:"about"`

	// Picture       []uint8 `json:"picture"`
}

type Session struct {
	User *User
	// NotificationQueue []*Comment
	NumNewMoves int
	//Blitz             *Match
	//Polling           bool
	//WsConn            *websocket.Conn
	//Data              chan interface{}
}

var SessionMap map[string]*Session // map of tokens to sessions

func Sessions() map[string]*Session {
	if SessionMap == nil {
		SessionMap = make(map[string]*Session)
	}
	return SessionMap
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

	return u
}
