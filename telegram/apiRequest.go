package telegram

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	ApiUrl = `https://api.telegram.org`
	RequestTimout = time.Second * 3
)

func MakeApiUrl(token, method string, parameters url.Values) (*url.URL, error) {
	apiUrl, err := url.Parse(fmt.Sprintf("%s/bot%s/%s", ApiUrl, token, method))
	if err != nil {
		log.Println(token, method, err.Error())
		return nil, err
	}
	apiUrl.RawQuery = parameters.Encode()
	return apiUrl, nil
}

func MakeGetApiRequest(token, method string, parameters url.Values) (int, []byte, error) {
	apiUrl, err := MakeApiUrl(token, method, parameters)
	if err != nil { return 0, nil, err }

	request, err := http.NewRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		log.Println(apiUrl.String(), err.Error())
		return 0, nil, err
	}

	client := http.Client{
		Timeout: RequestTimout,
	}
	response, err := client.Do(request)
	if err != nil {
		log.Println(apiUrl.String(), err.Error())
		return 0, nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(apiUrl.String(), err.Error())
		return 0, nil, err
	}

	return response.StatusCode, data, nil
}

func MakeGetApiRequestRetried(token, method string, parameters url.Values, retries int) (int, []byte, error) {
	retry := 0
	for {
		status, data, err := MakeGetApiRequest(token, method, parameters)
		if err == nil { return status, data, nil }
		retry = retry + 1
		if retry >= retries { break }
		log.Println(`retry`, retry, token, method, parameters)
	}
	msg := fmt.Sprintf("%d retries are fail!", retries)
	log.Println(msg, token, method, parameters)
	return 0, nil, errors.New(msg)
}