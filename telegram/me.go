package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Me struct {
	FirstName string `json:"first_name"`
	Id int `json:"id"`
	CanJoinGroups bool `json:"can_join_groups"`
	CanReadAllGroupMessages bool `json:"can_read_all_group_messages"`
	IsBot bool `json:"is_bot"`
	SupportsInlineQueries bool `json:"supports_inline_queries"`
	Username string `json:"username"`
}

func GetMe(token string) (Me, error) {
	status, data, err := MakeGetApiRequestRetried(token, `getMe`, url.Values{}, 3)
	if err != nil { return Me{}, err }
	var response struct {
		Ok bool `json:"ok"`
		Result Me `json:"result"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		log.Println(token, string(data), err.Error())
		return Me{}, err
	}

	if !response.Ok || status != http.StatusOK {
		msg := fmt.Sprintf("HTTP Status: %d; JSON ok: %v", status, response.Ok)
		log.Println(token, msg)
		return response.Result, errors.New(msg)
	}

	return response.Result, nil
}