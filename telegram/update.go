package telegram

import (
	"encoding/json"
	"log"
)

type Update struct {
	Id      int     `json:"update_id"`
	Message Message `json:"message"`
}

type Message struct {
	Id   int    `json:"message_id"`
	Date int    `json:"date"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
	Entities Entities `json:"entities"`
}

func (message *Message) IsBotCommand() bool {
	for _, v := range message.Entities {
		if v.Type == `bot_command` { return true }
	}
	return false
}

type Chat struct {
	Id                          int    `json:"id"`
	Type                        string `json:"type"`
	Title                       string `json:"title"`
	Username                    string `json:"username"`
	FirstName                   string `json:"first_name"`
	LastName                    string `json:"last_name"`
	AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
}

type Entities []struct {
	Type string `json:"type"`
	Offset int `json:"offset"`
	Length int `json:"length"`
}


func UnmarshalUpdate(data []byte) (Update, error) {
	update := Update{}
	if err := json.Unmarshal(data, &update); err != nil {
		log.Println(string(data), err.Error())
		return update, err
	}
	return update, nil
}
