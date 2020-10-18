package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
)

func DeleteMessage(token string, chatId, messageId int) error {
	parameters := url.Values{}
	parameters.Add(`chat_id`, fmt.Sprintf("%d", chatId))
	parameters.Add(`message_id`, fmt.Sprintf("%d", messageId))
	_, data, err := MakeGetApiRequestRetried(token, `deleteMessage`, parameters, 3)
	if err != nil { return err }
	var response struct{
		Ok bool `json:"ok"`
		ErrorCode int `json:"error_code"`
		Description string `json:"description"`
		Result bool `json:"result"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		log.Println(string(data), err.Error())
		return err
	}
	if !response.Ok {
		log.Println(string(data))
		return errors.New(response.Description)
	}
	return nil
}
