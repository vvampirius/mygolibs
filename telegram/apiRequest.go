package telegram

import (
	"fmt"
	"log"
	"net/url"
)

var (
	ApiUrl = `https://api.telegram.org`
)

func makeGetApiRequest(token, method string, parameters url.Values) ([]byte, error) {
	apiUrl, err := url.Parse(fmt.Sprintf("%s/bot%s/%s", ApiUrl, token, method))
	if err != nil {
		log.Println(token, method, err.Error())
		return nil, err
	}
	apiUrl.RawQuery = parameters.Encode()
	log.Println(API_URL)
	log.Println(apiUrl)
	log.Println(apiUrl.String())
	return nil, nil
}