package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type Api struct {
	Url string
	Token string
	RequestTimeout time.Duration
	RequestRetries int
	DebugLog *log.Logger
	ErrorLog *log.Logger
	ApiErrorFunc func(string, []byte)
}

func (api *Api) RequestUrl(method string) string {
	return fmt.Sprintf("%s/bot%s/%s", api.Url, api.Token, method)
}

func (api *Api) Do(method string, payload []byte) (int, []byte, error) {
	requestUrl := api.RequestUrl(method)
	buffer := bytes.NewBuffer(payload)

	request, err := http.NewRequest(http.MethodPost, requestUrl, buffer)
	if err != nil {
		if api.ErrorLog != nil { api.ErrorLog.Println(requestUrl, err.Error()) }
		return 0, nil, err
	}
	request.Header.Set(`Content-Type`, `application/json`)

	client := http.Client{
		Timeout: api.RequestTimeout,
	}
	response, err := client.Do(request)
	if err != nil {
		if api.ErrorLog != nil { log.Println(api.Url, method, string(payload), err.Error()) }
		if api.ApiErrorFunc != nil { api.ApiErrorFunc(method, payload) }
		return 0, nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		if api.ErrorLog != nil { log.Println(api.Url, method, string(payload), err.Error()) }
		return 0, nil, err
	}

	return response.StatusCode, data, nil
}

func (api *Api) DoWithRetry(method string, payload []byte) (int, []byte, error) {
	retry := 0
	for {
		status, data, err := api.Do(method, payload)
		if err != nil {
			if retry >= api.RequestRetries {
				time.Sleep(time.Second)
				return status, data, err
			}
			retry++
			continue
		}
		return status, data, err
	}
}

func (api *Api) Request(method string, payload interface{}) (int, *RequestResponse, error) {
	j, err := JsonEncode(payload)
	if err != nil {
		if api.ErrorLog != nil { log.Println(method, payload, err.Error()) }
		return 0, nil, err
	}
	status, data, err := api.DoWithRetry(method, j)
	if err != nil { return status, nil, err }
	requestResponse, err := api.UnmarshalRequestResponse(data)
	if err != nil { return status, nil, err}
	return status, requestResponse, err
}

func (api *Api) UnmarshalUpdate(data []byte) (*Update, error) {
	update := Update{}
	if err := json.Unmarshal(data, &update); err != nil {
		if api.ErrorLog != nil { api.ErrorLog.Println(string(data), err.Error()) }
		return nil, err
	}
	return &update, nil
}

func (api *Api) UnmarshalRequestResponse(data []byte) (*RequestResponse, error) {
	requestResponse := RequestResponse{}
	if err := json.Unmarshal(data, &requestResponse); err != nil {
		if api.ErrorLog != nil { api.ErrorLog.Println(string(data), err.Error()) }
		return nil, err
	}
	return &requestResponse, nil
}

func (api *Api) SendFile(method, fileName string, r io.Reader, fields map[string]string) (int, *RequestResponse, error) {
	tempFile, err := ioutil.TempFile(os.TempDir(), fileName)
	if err != nil {
		if api.ErrorLog != nil { api.ErrorLog.Println(err.Error()) }
		return 0, nil, err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	if api.DebugLog != nil { api.DebugLog.Printf("Temp file '%s' created", tempFile.Name()) }
	//defer os.Remove(tempFile.Name())
	writer := multipart.NewWriter(tempFile)
	for field, value := range fields {
		if err := writer.WriteField(field, value); err != nil {
			if api.ErrorLog != nil { api.ErrorLog.Println(err.Error()) }
			return 0, nil, err
		}
	}
	filePart, err := writer.CreateFormFile(`file`, fileName)
	if err != nil {
		if api.ErrorLog != nil { api.ErrorLog.Println(err.Error()) }
		return 0, nil, err
	}
	if _, err := io.Copy(filePart, r); err != nil {
		if api.ErrorLog != nil { api.ErrorLog.Println(err.Error()) }
		return 0, nil, err
	}
	if err := writer.Close(); err != nil {
		if api.ErrorLog != nil { api.ErrorLog.Println(err.Error()) }
		return 0, nil, err
	}

	tempFile.Seek(0, 0)
	requestUrl := api.RequestUrl(method)
	request, err := http.NewRequest(http.MethodPost, requestUrl, tempFile)
	if err != nil {
		if api.ErrorLog != nil { api.ErrorLog.Println(requestUrl, err.Error()) }
		return 0, nil, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := http.Client{
		Timeout: api.RequestTimeout,
	}
	response, err := client.Do(request)
	if err != nil {
		if api.ErrorLog != nil { log.Println(api.Url, method, tempFile.Name(), err.Error()) }
		//if api.ApiErrorFunc != nil { api.ApiErrorFunc(method, tempFile.Name()) }
		return 0, nil, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		if api.ErrorLog != nil { log.Println(api.Url, method, tempFile.Name(), err.Error()) }
		return 0, nil, err
	}
	requestResponse, err := api.UnmarshalRequestResponse(responseBody)
	if err != nil {
		if api.ErrorLog != nil { log.Println(api.Url, method, tempFile.Name(), err.Error()) }
		return 0, nil, err
	}

	return response.StatusCode, requestResponse, nil
}

func NewApi(token string) *Api {
	api := Api{
		Url: DefaultApiUrl,
		Token: token,
		RequestTimeout: DefaultRequestTimeout,
		RequestRetries: DefaultRequestRetries,
	}
	return &api
}