package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func SetWebHook(token, callbackUrl string) error {
	parameters := url.Values{}
	parameters.Add(`url`, callbackUrl)
	status, data, err := MakeGetApiRequestRetried(token, `setWebHook`, parameters, 3)
	if err != nil { return err }

	var response struct {
		Ok bool `json:"ok"`
		Result bool `json:"result"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		log.Println(token, string(data), err.Error())
		return err
	}

	if !response.Ok || status != http.StatusOK {
		msg := fmt.Sprintf("HTTP Status: %d; JSON ok: %v; Description: %s", status, response.Ok, response.Description)
		log.Println(token, msg)
		return errors.New(msg)
	}

	return nil
}
