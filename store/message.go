package store

import (
	"github.com/jarrancarr/ChexxServer/utils"
	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	Meta       string `json:"meta"` // room, type, urgency, related... etc
	Topic      string `json:"topic"`
	Body       string `json:"body"`
	Author     uint   `json:"author"`
	Recipients string `json:"recipients"`
	CC         string `json:"cc"`
	BCC        string `json:"bcc"`
	Attach     string `json:"attach"`
}

func (m *Message) Create() map[string]interface{} {
	GetDB().Create(m)
	if m.ID <= 0 {
		return utils.Message(false, "Failed to create message, connection error.")
	}
	response := utils.Message(true, "Message created")
	response["message"] = m
	return response
}
func GetMessage(id uint) *Message {

	m := &Message{}
	GetDB().Table("messages").Where("id = ?", id).First(m)
	// json.Unmarshal([]byte(u.Property), &u.Prop)
	// if u.Prop == nil {
	// 	u.Prop = make(map[string]string)
	// 	u.Prop["test"] = "success"
	// }
	return m
}
func (m *Message) Validate() (map[string]interface{}, bool) {

	if m.Author == 0 {
		return utils.Message(false, "No Origin"), false
	}

	return utils.Message(true, "OK"), true
}
func (m *Message) Update() map[string]interface{} {

	if resp, ok := m.Validate(); !ok {
		return resp
	}

	// convert properties

	// prop, err := json.Marshal(u.Prop)
	// if err == nil {
	// 	u.Property = string(prop)
	// }

	GetDB().Save(m)

	if m.ID <= 0 {
		return utils.Message(false, "Failed to create message, connection error.")
	}

	response := utils.Message(true, "Message updated")
	return response
}
func (m *Message) Delete() map[string]interface{} {
	GetDB().Delete(m)
	response := utils.Message(true, "Message deleted")
	return response
}
