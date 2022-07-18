package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/jarrancarr/ChexxServer/store"
	"github.com/jarrancarr/ChexxServer/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var DEBUG = true

func FindUser(r *http.Request) (*store.User, error) {
	if DEBUG {
		log.Println("FindUser")
	}
	token := r.Header.Get("Authorization")
	// fmt.Println(">>>" + token + "<<<")
	if store.Sessions()[token] == nil {
		return nil, errors.New("No User by that ID found")
	}
	user := store.Sessions()[token].User
	user.Token = token
	return user, nil
}

func GetUser(u uint) *store.User {
	if DEBUG {
		log.Println("GetUser")
	}
	user := &store.User{}
	store.GetDB().Table("users").Where("id = ?", u).First(user)
	if user.Email == "" { //User not found!
		return nil
	}

	user.Prop = make(map[string]string)

	// convertProps(user)

	user.Password = ""
	return user
}

func GetUserByUserIdOrEmail(uid, email string) *store.User {
	if DEBUG {
		log.Println("GetUserByUserIdOrEmail")
	}
	user := &store.User{}

	err := store.GetDB().Table("users").Where("user_id = ?", uid).First(user).Error
	if err != nil {
		err = store.GetDB().Table("users").Where("email = ?", email).First(user).Error
		if err != nil {
			return nil
		}
	}

	user.Prop = make(map[string]string)
	// convertProps(user)
	user.Password = ""
	return user
}

var Facebook = func(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("Facebook")
	}
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data) //decode the request body into struct and failed if any error occur
	if err != nil {
		utils.Respond(w, utils.Message(false, "Problem decoding request"))
		return
	}

	user := GetUserByUserIdOrEmail("", fmt.Sprintf("%v", data["email"]))
	user.Token = fmt.Sprintf("%v", data["accessToken"])
	fmt.Printf("User %s facebook logged in\n", user.Name)
	store.Sessions()[user.Token] = &store.Session{User: user, NumNewMoves: 0}

	resp := utils.Message(true, "Logged In")
	resp["user"] = user

	utils.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("Authenticate")
	}
	user := &store.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}

	resp := Login(user.UserId, user.Password)
	utils.Respond(w, resp)
	// if resp["status"] == true {
	// 	purgeBlitzGames(resp["account"].(*models.Account).ID)
	// }
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("Logout")
	}
	user, err := FindUser(r)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}
	delete(store.Sessions(), user.Token)
	delete(store.Online(), user.ID)
	utils.Respond(w, utils.Message(true, "Logged Out"))
	fmt.Printf("user %s logged out\n", user.Name)
}
func Login(userId, password string) map[string]interface{} {
	if DEBUG {
		log.Println("Login")
	}

	// is user already logged in with a session?

	user := &store.User{}
	err := store.GetDB().Table("users").Where("user_id = ?", userId).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Message(false, "User Id address not found")
		}
		return utils.Message(false, "Connection error. Please retry")
	}

	// fmt.Println(user)
	user.Revert()
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return utils.Message(false, "Invalid login credentials. Please try again")
	}
	user.Password = ""
	if DEBUG {
		log.Println("...Password checked")
	}
	//Create JWT token
	tk := &store.Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString //Store the token in the response

	// convertProps(account)

	// remove old login stuff

	oldToken := store.Online()[user.ID]
	if oldSession, ok := store.Sessions()[oldToken]; ok {
		oldSession.Inbox <- true
		if DEBUG {
			log.Println("...old session removed")
		}
	} // kill wsDataQueue
	delete(store.Sessions(), oldToken)
	delete(store.Online(), user.ID)
	delete(store.BlitzMatches(), oldToken)
	if DEBUG {
		log.Println("...old blitz stuff removed")
		log.Printf("online:\n   %v\n", store.Online())
		log.Printf("sessions:\n   %v\n", store.Sessions())
		log.Printf("BlitzMatches()\n   %v\n", store.BlitzMatches())
	}
	store.Sessions()[user.Token] = &store.Session{User: user, NumNewMoves: 0, Inbox: make(chan interface{})}
	store.Online()[user.ID] = user.Token

	resp := utils.Message(true, "Logged In")
	resp["user"] = user
	if DEBUG {
		log.Println("...completed")
	}
	return resp
}
func UserInfo(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("UserInfo")
	}
	user, _ := FindUser(r)

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Respond(w, utils.Message(false, "ID not found."))
		return
	}
	u := store.GetUser(uint(id))

	resp := utils.Message(true, "Found match")
	resp["opponent"] = u

	utils.Respond(w, resp)
}
func Queues(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("Queues")
	}
	user, _ := FindUser(r)

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}
	messages := GetMessagesFor(user.UserId)
	requests := GetFriendRequestsFor(user.UserId)
	resp := utils.Message(true, "Queues")
	resp["messages"] = convert(messages)
	resp["requests"] = convert(requests)

	utils.Respond(w, resp)
}
func convert(msgs []store.Message) []interface{} {
	if DEBUG {
		log.Println("convert")
	}
	data := make([]interface{}, len(msgs))
	for m := range msgs {
		author := GetUser(msgs[m].Author)
		data[m] = struct {
			ID    uint   `json:"ID"`
			Meta  string `json:"meta"`
			From  string `json:"from"`
			Topic string `json:"topic"`
			Text  string `json:"text"`
		}{ID: msgs[m].ID, Meta: msgs[m].Meta, From: author.Name, Topic: msgs[m].Topic, Text: msgs[m].Body}
	}
	return data
}
func GetFriendRequestsFor(uId string) []store.Message {
	if DEBUG {
		log.Println("GetFriendRequestsFor")
	}
	requests := []store.Message{}
	err := store.GetDB().Table("messages").Where("meta = ? and recipients = ?", "FriendRequest", uId).Find(&requests).Error
	if err != nil {
		return nil
	}
	return requests
}
func GetMessagesFor(uId string) []store.Message {
	if DEBUG {
		log.Println("GetMessagesFor")
	}
	messages := []store.Message{}
	err := store.GetDB().Table("messages").Where("recipients LIKE ?", "%"+uId+"%").Find(&messages).Error
	if err != nil {
		return nil
	}
	return messages
}
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("RegisterUser")
	}
	user := &store.User{}
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request!! "+err.Error()))
		return
	}
	occupied := GetUserByUserIdOrEmail(user.UserId, user.Email)
	if occupied != nil {
		utils.Respond(w, utils.Message(false, "userId or email address already used."))
		return
	}
	password := user.Password
	resp := user.Create()

	if resp["status"] == true {
		resp := Login(user.UserId, password)
		utils.Respond(w, resp)
	} else {
		utils.Respond(w, resp)
	}
}

var JwtAuthentication = func(next http.Handler) http.Handler {
	if DEBUG {
		log.Println("JwtAuthentication")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/tutor", "/user", "/match/cpu", "/pub", "/ws"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                                           //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if strings.HasPrefix(requestPath, value) {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			response = utils.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		// if len(splitted) != 2 {
		// 	response = utils.Message(false, "Invalid/Malformed auth token")
		// 	w.WriteHeader(http.StatusForbidden)
		// 	w.Header().Add("Content-Type", "application/json")
		// 	utils.Respond(w, response)
		// 	return
		// }

		tokenPart := splitted[0] //Grab the token part, what we are truly interested in
		tk := &store.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = utils.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			response = utils.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Printf("User %v/n", tk.UserId) //Useful for monitoring
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

func FriendAccept(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("FriendAccept")
	}
	data := &struct {
		Id int `json:"requestId"`
	}{}
	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body for new friend request"))
		return
	}

	user, _ := FindUser(r)

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}
	request := store.GetMessage(uint(data.Id))

	for req := range user.Friend {
		if user.Friend[req].ID == request.Author {
			utils.Respond(w, utils.Message(false, "Already Friends"))
			request.Delete()
			return
		}
	}

	requestor := GetUser(request.Author)
	if requestor != nil {
		newFriend := store.Friend{ID: requestor.ID, Name: requestor.Name, UserId: requestor.UserId, Colors: "#8bb|#466|#fda"}
		beaFriend := store.Friend{ID: user.ID, Name: user.Name, UserId: user.UserId, Colors: "#8bb|#466|#fda"}

		user.Friend = append(user.Friend, newFriend)
		requestor.Friend = append(requestor.Friend, beaFriend)
		user.Update()
		requestor.Update()
		resp := utils.Message(true, "You have a new friend.")
		resp["friend"] = newFriend
		request.Delete()
		utils.Respond(w, resp)
	} else {
		utils.Respond(w, utils.Message(false, "user not found."))
	}
}
func FriendReject(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("FriendReject")
	}
	data := &struct {
		Id int `json:"requestId"`
	}{}
	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body for friend rejection"))
		return
	}

	user, _ := FindUser(r)

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}
	request := store.GetMessage(uint(data.Id))
	resp := request.Delete()
	utils.Respond(w, resp)
}
func MessageOk(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("MessageOk")
	}
	user, _ := FindUser(r)

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Respond(w, utils.Message(false, "ID not found."))
		return
	}
	request := store.GetMessage(uint(id))
	resp := request.Delete()
	utils.Respond(w, resp)
}
func FriendRequest(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("FriendRequest")
	}
	var jsonData map[string]string
	err := json.NewDecoder(r.Body).Decode(&jsonData)

	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body for new friend request"))
		return
	}

	user, _ := FindUser(r)

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}
	recipient := GetUserByUserIdOrEmail(jsonData["userid"], "")
	fr := GetFriendRequestsFor(recipient.UserId)
	for req := range fr {
		if fr[req].Recipients == recipient.UserId {
			utils.Respond(w, utils.Message(false, "Duplicate Friend Request"))
			return
		}
	}
	newMessage := store.Message{Meta: "FriendRequest", Author: user.ID, Body: jsonData["message"], Recipients: recipient.UserId}

	newMessage.Create()

	if session, ok := store.Sessions()[store.OnlineMapping[recipient.ID]]; ok {
		session.Inbox <- newMessage
	}

	utils.Respond(w, utils.Message(true, "OK"))
}
func Message(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("Message")
	}
	msg := &struct {
		Meta  string `json:"meta"`
		Topic string `json:"topic"`
		Text  string `json:"body"`
		To    string `json:"recipient"`
	}{}

	err := json.NewDecoder(r.Body).Decode(msg)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body for new comment"))
		return
	}

	user, _ := FindUser(r)

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}
	recipient := GetUserByUserIdOrEmail(msg.To, "")
	newMessage := store.Message{Meta: msg.Meta, Author: user.ID, Topic: msg.Topic, Body: msg.Text, Recipients: msg.To}

	newMessage.Create()

	if session, ok := store.Sessions()[store.OnlineMapping[recipient.ID]]; ok {
		session.Inbox <- newMessage
	}

	utils.Respond(w, utils.Message(true, "OK"))
}
func SaveUser(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Println("SaveUser")
	}
	modUser := &store.User{}
	err := json.NewDecoder(r.Body).Decode(modUser)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body for new comment"))
		return
	}

	user, _ := FindUser(r)

	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found."))
		return
	}

	user.Prop = modUser.Prop
	resp := user.Update()
	utils.Respond(w, resp)
}
