package store

import (
	"fmt"
	"math/rand"

	"github.com/jinzhu/gorm"
)

//a struct to rep user account
type Test struct {
	gorm.Model
	Who   string `json:"who"`
	What  string `json:"what"`
	Where string `json:"where"`
	When  string `json:"when"`
	Why   string `json:"why"`
	How   string `json:"how"`
}

func (t *Test) Create() {

	GetDB().Create(t)

	if t.ID <= 0 {
		fmt.Println("Failed to create test, connection error.")
	}

	fmt.Println("Test has been created")
}

func GetTest(u uint) *Test {

	t := &Test{}
	GetDB().Table("tests").Where("id = ?", u).First(t)
	if t.Who == "" || t.What == "" || t.Where == "" { //User not found!
		return nil
	}

	return t
}

func DoTest() {
	subject := Test{Who: Gibber([]int{9, 12}), What: Scrabble(8), Where: Gibber([]int{5, 2, 9, 5}), When: Gibber([]int{7, 4}), Why: Gibber([]int{6, 7, 3, 6, 9, 4, 3, 4, 7}), How: Gibber([]int{9, 5, 3, 3, 2, 9, 6})}
	subject.Create()

	object := GetTest(subject.ID)

	fmt.Printf("subject %d::%v\n", subject.ID, subject)
	fmt.Printf("object %d::%v\n", object.ID, object)
}

func Gibber(n []int) string {
	sentence := ""
	for x := range n {
		sentence += Scrabble(x) + " "
	}
	sentence = sentence[0:len(sentence)-1] + "."
	return sentence
}
func Scrabble(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZaaaabbcccddeeeeefghhhiiiijkkllmmmnnnnooooopppqrrrrssssttttuuuvwwxyyyz")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
