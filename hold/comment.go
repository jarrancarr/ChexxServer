package hold

import (
	"fmt"
	"math"
	"strings"

	u "Chexx/server/utils"

	"github.com/jinzhu/gorm"
)

var logBuffer = ""
var INDENT = 0
var DEBUG = 5 // 0-Off  1-Error  2-Info  3-Warn  4-Debug  5-Trace 6-GodHelpYou
//  log(X, fmt.Sprintf("<function(%s):returnType>  ", params...))

func log(level int, log string) {
	buffer := false
	if level < 0 {
		buffer = true
		level = -level
	}
	if level <= DEBUG {
		if strings.Index(log, "</") > -1 {
			INDENT--
		}
		if INDENT > -1 {
			if buffer {
				if log == "" {
					fmt.Println("|  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  "[:INDENT*3] + logBuffer)
					logBuffer = ""
				} else {
					logBuffer += log
				}
			} else {
				fmt.Println("|  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  "[:INDENT*3] + logBuffer + log)
				logBuffer = ""
			}
		}
		if strings.Index(log, "</") == -1 && strings.Index(log, "/>") == -1 && strings.Index(log, ">") > 0 && strings.Index(log, "<") > -1 && strings.Index(log, ">") > strings.Index(log, "<") {
			INDENT++
		}
	}
	if INDENT > 20 {
		INDENT = 20
	}
	if INDENT < 0 {
		INDENT = 0
	}
}

type COMMENTTYPE int

const (
	GENERAL = iota
	REPLY
	CHATROOM
	DIRECT
	MATCH
	LEADER
	MEMBER
	INVITE
	NOTIFICATION
	NOTE_READ
	NOTE_ACKNOWLEDGED
	NOTE_REPLIED
	NOTE_ARCHIVED
)

type Comment struct {
	gorm.Model
	RefType uint
	RefId   uint
	Topic   string
	Tags    string
	Content string
	Author  uint
}

type Message struct {
	gorm.Model
	Room    string `json:"name"`
	Message string `json:"message"`
	Author  uint   `json:"author"`
}

type Room struct {
	Queue    map[*Message]int64 // if divisible by prime, message received: so 455 means= Role ids  [5, 7, 13] received message
	Role     map[uint]int64     // AccountId : prime designator, each active user is assigned a prime number to track messages received.
	Overflow *Room              // primes up to 47 are realistic.  So if more, create an adjoining room
	Topics   map[string]int
}

func (r *Room) Announce(m *Message) {
	me, ok := r.Role[m.Author]
	if ok {
		r.Queue[m] = me // initialize this message as marked read by the author by his prime designator
	} else {
		r.Queue[m] = 1 // user not in this room, so set prime designator to one so all get it.
	}
	// iterate through members that have sessions.
	for id, prm := range r.Role {
		if id != m.Author {
			ActiveSession[id].Data <- m
			r.Queue[m] *= prm
		}
	}
	if r.Overflow != nil {
		r.Announce(m)
	}
}

func (r *Room) Hashtag(m *Message) {
	debug := 5
	log(debug, fmt.Sprintf("<Room::Hashtag>"))
	log(debug, fmt.Sprintf("message: %v", m))
	for _, w := range strings.Split(m.Message, " ") {
		log(debug, fmt.Sprintf("checking %s", w))
		tokens := strings.Split(w, "|||")
		for _, tc := range tokens {
			log(debug, fmt.Sprintf("  token %s", tc))
			if strings.HasPrefix(tc, "#") {
				log(debug, fmt.Sprintf("    is a hashtag?"))
				found := false
				for t, pts := range r.Topics {
					log(debug, fmt.Sprintf("      ?= %s", t))
					if t == tc {
						log(debug, fmt.Sprintf("        yes"))
						found = true
						r.Topics[t] = pts + 1
					}
				}
				log(debug, fmt.Sprintf("    "))
				if !found {
					log(debug, fmt.Sprintf("        no"))
					r.Topics[tc] = 1
					for id, _ := range r.Role {
						ActiveSession[id].Data <- "hashtag|||" + m.Room + "##" + tc
					}
				}
			}
		}
	}
	log(debug, fmt.Sprintf("</Room::Hashtag>  "))
}

func (r *Room) Enter(id uint) int64 {
	good := false
	var me int64
	for me = 1; !good; {
		if me > 47 { // already too many in this room
			if r.Overflow == nil {
				r.Overflow = &Room{Queue: make(map[*Message]int64), Role: make(map[uint]int64), Overflow: nil, Topics: make(map[string]int)}
			}
			return 0
		}
		me += 2
		good = true
		for factor := int64(3); float64(factor) <= math.Sqrt(float64(me)) && good; factor += 2 {
			good = me%factor != 0
		}
		for _, pm := range r.Role {
			if me == pm {
				good = false // someone already has that number
			}
		}
	}
	r.Role[id] = me
	return me
}

func (r *Room) GetChats(uid uint) []Message {
	if me, ok := r.FindUser(uid); ok {
		messages := make([]Message, 0)
		for m, v := range r.Queue {
			if v%me > 0 {
				r.Queue[m] = v * me // mark message as read
				messages = append(messages, *m)
			}
		}
		return messages
	}
	return nil
}

func (r *Room) FindUser(id uint) (int64, bool) {
	me, ok := r.Role[id]
	if !ok && r.Overflow != nil {
		return r.Overflow.FindUser(id)
	}
	return me, ok
}

func (c Comment) String() string {
	return fmt.Sprintf("%s // %s // %d : %s", c.Topic, c.Tags, c.Author, c.Content)
}

func (c Message) String() string {
	return fmt.Sprintf("%s // %d : %s", c.Room, c.Author, c.Message)
}

//Validate incoming user details...
func (c *Comment) Validate() (map[string]interface{}, bool) {

	if c.Author < 1 {
		return u.Message(false, "No Author"), false
	}

	return u.Message(false, "Requirement passed"), true
}

//Validate incoming user details...
func (c *Message) Validate() (map[string]interface{}, bool) {

	if c.Author < 1 {
		return u.Message(false, "No Author"), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (c *Comment) Create() map[string]interface{} {

	if resp, ok := c.Validate(); !ok {
		return resp
	}

	GetDB().Create(c)

	if c.ID <= 0 {
		return u.Message(false, "Failed to create comment, connection error.")
	}

	response := u.Message(true, "Comment created")
	response["id"] = c.ID
	return response
}

func (c *Message) Create() map[string]interface{} {

	if resp, ok := c.Validate(); !ok {
		return resp
	}

	GetDB().Create(c)

	if c.ID <= 0 {
		return u.Message(false, "Failed to create message, connection error.")
	}

	response := u.Message(true, "Comment created")
	response["id"] = c.ID
	return response
}

func (c *Comment) Update() map[string]interface{} {
	GetDB().Save(c)
	response := u.Message(true, "Comment has been updated")
	return response
}

func GetComment(u uint) *Comment {

	com := &Comment{}
	GetDB().Table("comments").Where("id = ?", u).First(com)
	return com
}

func GetMessage(u uint) *Message {

	com := &Message{}
	GetDB().Table("messages").Where("id = ?", u).First(com)
	return com
}

func GetMessageLog(root string) []Message {

	logs := []Message{}
	GetDB().Table("messages").Where("room = ?", root).Find(&logs)
	return logs
}

func GetCommentsByReference(rid uint, depth int) []Comment {
	return GetCommentsByReferences([]uint{rid}, depth)
}

func GetCommentsByReferences(rid []uint, depth int) []Comment {
	var coms []Comment
	GetDB().Table("comments").Where("ref_id IN (?)", rid).Find(&coms)

	if depth > 0 {
		var subComs []Comment
		srid := make([]uint, len(coms))
		for _, c := range subComs {
			srid = append(srid, c.RefId)
		}
		subComs = GetCommentsByReferences(srid, depth-1)
		return append(coms, subComs...)
	}

	return coms
}

func GetCommentsByTopic(topic string) []Comment {

	var coms []Comment

	GetDB().Table("comments").Where("topic = ?", topic).Find(&coms)

	return coms
}

var Notification = func(topic, message string, userId uint) uint {
	note := Comment{Author: userId, RefType: NOTIFICATION, Content: message, Topic: topic}
	GetDB().Create(&note)
	ActiveSession[userId].Notify(&note)
	return note.ID
}

func Tag(noteId uint, tag string) {
	note := GetComment(noteId)
	note.Tags = note.Tags + "!!!" + tag
	note.Update()
}

func GetCommentsByType(userId uint, refType COMMENTTYPE) []Comment {
	var coms []Comment
	GetDB().Table("comments").Where("REF_TYPE = ? AND Author = ?", refType, userId).Find(&coms)
	return coms
}

func GetAllComments() []Comment {
	var coms []Comment
	GetDB().Table("comments").Where("REF_TYPE = 0").Find(&coms)
	return coms
}
