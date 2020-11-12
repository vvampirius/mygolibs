package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

type SendMessageError struct {
	ErrCode int
	Description string
}

func (sendMessageError *SendMessageError) Error() string {
	return sendMessageError.Description
}

func SendMessage(token string, chatId int, text string, disableNotification bool, replyToMessageId int) (Message, error) {
	parameters := url.Values{}
	parameters.Add(`chat_id`, fmt.Sprintf("%d", chatId))
	parameters.Add(`text`, text)
	if disableNotification { parameters.Add(`disable_notification`, `True`)}
	if replyToMessageId > 0 { parameters.Add(`reply_to_message_id`, fmt.Sprintf("%d", replyToMessageId)) }
	_, data, err := MakeGetApiRequestRetried(token, `sendMessage`, parameters, 3)
	if err != nil { return Message{}, err }
	var response struct{
		Ok bool `json:"ok"`
		Description string `json:"description"`
		ErrorCode int `json:"error_code"`
		Result Message `json:"result"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		log.Println(string(data), err.Error())
		return Message{}, err
	}
	if !response.Ok {
		err := SendMessageError{
			ErrCode: response.ErrorCode,
			Description: response.Description,
		}
		log.Printf("Send message to chats %d got error: %s\n", chatId, string(data))
		return Message{}, &err
	}
	return response.Result, nil
}
