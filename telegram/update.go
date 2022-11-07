package telegram

import (
	"encoding/json"
	"log"
)

type Update struct {
	Id      int     `json:"update_id"`
	Message Message `json:"message"`
	CallbackQuery CallbackQuery `json:"callback_query"`
}

func (update *Update) IsCallbackQuery() bool {
	if update.CallbackQuery.Id != `` { return true }
	return false
}

func (update *Update) IsMessage() bool {
	if update.Message.Id != 0 { return true }
	return false
}

type Message struct {
	Id   int    `json:"message_id"`
	From User	`json:"from"`
	Date int    `json:"date"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
	Entities Entities `json:"entities"`
	ReplyToMessage ReplyToMessage `json:"reply_to_message"`
}

func (message *Message) IsBotCommand() bool {
	for _, v := range message.Entities {
		if v.Type == `bot_command` { return true }
	}
	return false
}

func (message *Message) IsReplyToMessage() bool {
	if message.ReplyToMessage.Id != 0 && message.ReplyToMessage.Date != 0 { return true }
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

type CallbackQuery struct {
	ChatInstance string `json:"chat_instance"`
	Data string `json:"data"`
	Id string `json:"id"`
	Message Message `json:"message"`
}

type User struct {
	Id int				`json:"id"`
	IsBot bool			`json:"is_bot"`
	FirstName string	`json:"first_name"`
	LastName string		`json:"last_name"`
	Username string 	`json:"username"`
	LanguageCode string	`json:"language_code"`
}

type ReplyToMessage struct {
	Id   int    `json:"message_id"`
	From User	`json:"from"`
	Date int    `json:"date"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
	Entities Entities `json:"entities"`
}

func UnmarshalUpdate(data []byte) (Update, error) {
	update := Update{}
	if err := json.Unmarshal(data, &update); err != nil {
		log.Println(string(data), err.Error())
		return update, err
	}
	return update, nil
}
