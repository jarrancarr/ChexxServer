package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/jarrancarr/ChexxServer/store"
	"github.com/jarrancarr/ChexxServer/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func FindUser(r *http.Request) (*store.User, error) {
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data) //decode the request body into struct and failed if any error occur
	if err != nil {
		return nil, errors.New("problem getting token from json")
	}
	token := fmt.Sprintf("%v", data["token"])
	fmt.Println(">>>" + token + "<<<")
	if store.Sessions()[token] == nil {
		return nil, errors.New("No User by that ID found")
	}
	user := store.Sessions()[token].User
	return user, nil
}

func GetUser(u uint) *store.User {

	user := &store.User{}
	store.GetDB().Table("users").Where("id = ?", u).First(user)
	if user.Email == "" { //User not found!
		return nil
	}

	user.Property = make(map[string]string)

	// convertProps(user)

	user.Password = ""
	return user
}

func GetUserByUserIdOrEmail(uid, email string) *store.User {

	user := &store.User{}

	err := store.GetDB().Table("users").Where("user_id = ?", uid).First(user).Error
	if err != nil {
		err = store.GetDB().Table("users").Where("email = ?", email).First(user).Error
		if err != nil {
			return nil
		}
	}

	user.Property = make(map[string]string)
	// convertProps(user)
	user.Password = ""
	return user
}

var Facebook = func(w http.ResponseWriter, r *http.Request) {

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
	fmt.Println("logging out")
	user, err := FindUser(r)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}
	store.Sessions()[user.Token] = nil
	utils.Respond(w, utils.Message(true, "Logged Out"))
	fmt.Printf("user %s logged out\n", user.Name)
}

func Login(userId, password string) map[string]interface{} {

	// is user already logged in with a session?

	user := &store.User{}
	err := store.GetDB().Table("users").Where("user_id = ?", userId).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Message(false, "User Id address not found")
		}
		return utils.Message(false, "Connection error. Please retry")
	}

	fmt.Println(user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return utils.Message(false, "Invalid login credentials. Please try again")
	}
	user.Password = ""

	//Create JWT token
	tk := &store.Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString //Store the token in the response

	// convertProps(account)

	fmt.Printf("User %s logged in\n", user.Name)
	store.Sessions()[user.Token] = &store.Session{User: user, NumNewMoves: 0}

	resp := utils.Message(true, "Logged In")
	resp["user"] = user
	return resp
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {

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

	fmt.Println("-------------RegisterUser-------------")
	fmt.Printf("%v\n", user.UserId)
	fmt.Printf("%v\n", password)
	fmt.Println("-------------RegisterUser-------------")

	if resp["status"] == true {
		resp := Login(user.UserId, password)
		utils.Respond(w, resp)
	} else {
		utils.Respond(w, resp)
	}
}

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/tutor", "/user"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path              //current request path

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