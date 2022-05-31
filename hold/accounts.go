package hold

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	u "Chexx/server/utils"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var LoginCount = 0

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

type Roles int

const (
	ADMIN = iota
	ANONYMOUS
	PREMIUM
	BENEFACTOR
)

//a struct to rep user account
type Account struct {
	gorm.Model
	Email    string            `json:"email"`
	Name     string            `json:"name"`
	Title    string            `json:"title"`
	Password string            `json:"password"`
	Ranks    string            `json:"ranks"`
	Token    string            `json:"token" gorm:"-"`
	Game     map[string]*Match `gorm:"-"`
	Lord     uint              `json:"lord"`
	About    string            `json:"about"`
	Pict     string            `json:"pict"`
	Role     Roles             `json:"role"`
	School   string            `json:"school"`
	City     string            `json:"city"`
	Company  string            `json:"company"`
	DOB      uint              `json:"dob"`
	Origin   string            `json:"origin"`
	Profile  string            `json:"profile"`
	Property map[string]string `json:"property" gorm:"-"`
	Props    string
}

type Session struct {
	User              *Account
	NotificationQueue []*Comment
	NumNewMoves       int
	Blitz             *Match
	Polling           bool
	WsConn            *websocket.Conn
	Data              chan interface{}
}

type Team struct {
	gorm.Model
	Name      string              `json:"name"`
	Officers  map[string]*Account `gorm:"-"`
	Members   map[string]*Account `gorm:"-"`
	Applicant map[string]*Account `gorm:"-"`
	About     string              `json:"about"`
	Pict      string              `json:"pict"`
	Rule      map[string]string   `gorm:"-"`
	Room      *Room               `gorm:"-"`
}

var ActiveSession map[uint]*Session

func (s *Session) Notify(note *Comment) {
	if s == nil {
		fmt.Println("session is null")
		return
	}
	if s.NotificationQueue == nil {
		fmt.Println("NotificationQueue is null")
		return
	}
	s.NotificationQueue = append(s.NotificationQueue, note)
}

func (s *Session) PollNotifications() []*Comment {
	pull := s.NotificationQueue
	s.NotificationQueue = make([]*Comment, 0)
	return pull
}

func (a Account) String() string {
	return fmt.Sprintf("%s // %s // %s", a.Name, a.Title, a.Email)
}

//Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	account.Ranks = "1000"
	account.Role = ANONYMOUS

	for i := 1; i < LASTTYPE*LASTFORMATION; i++ {
		account.Ranks += "|1000"
	}

	serializeProps(account)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	//Create new JWT token for the newly registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func Login(email, password string) map[string]interface{} {
	debug := 6
	log(debug, fmt.Sprintf("<Login> %s::%s", email, password))

	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	log(debug, fmt.Sprintf("found account: password=%s", account.Password))

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		log(debug, fmt.Sprintf("</Login> no good, bad password"))
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	LoginCount += 1
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	if ActiveSession == nil {
		ActiveSession = make(map[uint]*Session)
	}

	convertProps(account)

	ActiveSession[account.ID] = &Session{User: account, NotificationQueue: make([]*Comment, 0), NumNewMoves: 0}

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	log(debug, fmt.Sprintf("</Login>"))

	return resp
}

func CountUsers() int {
	var count int
	GetDB().Model(&Account{}).Count(&count)
	return count
}

func GetUser(u uint) *Account {

	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Property = make(map[string]string)

	convertProps(acc)

	acc.Password = ""
	return acc
}

func FindUser(r *http.Request) (*Account, error) {
	userId := r.Header.Get("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		return nil, err
	}
	user := GetUser(uint(id))
	if user == nil {
		return nil, errors.New("No User by that ID found")
	}
	return user, nil
}

func (account *Account) Update() map[string]interface{} {
	serializeProps(account)
	if account.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
		account.Password = string(hashedPassword)
	}
	GetDB().Save(account)
	response := u.Message(true, "account has been updated")
	response["account"] = account
	return response
}

func (user *Account) GetRating(match *Match) int {
	debug := 6
	log(debug, fmt.Sprintf("<GetRating> %v", match))
	log(debug, fmt.Sprintf("user ratings: %s", user.Ranks))
	ranks := strings.Split(user.Ranks, "|")
	if len(ranks) > int(match.GType)*int(LASTFORMATION)+int(match.Formation) {
		rnk, err := strconv.Atoi(ranks[int(match.GType)*int(LASTFORMATION)+int(match.Formation)])
		if err != nil {
			log(debug, "</GetRating> Error")
			return -1
		}
		log(debug, "</GetRating>")
		return rnk
	}
	log(debug, "</GetRating> Problem")
	return 1000
}

func Announce(json string) {
	for acc := range ActiveSession {
		ActiveSession[acc].Data <- json
	}
}

func convertProps(acc *Account) {
	acc.Property = make(map[string]string)
	for _, r := range strings.Split(acc.Props, " ") {
		kv := strings.Split(r, ":")
		if len(kv) == 2 {
			acc.Property[kv[0]] = kv[1]
		}
	}
	acc.Props = ""
}

func serializeProps(acc *Account) {
	acc.Props = ""
	for k, p := range acc.Property {
		acc.Props += k + ":" + p + " "
	}
	acc.Property = nil
}
