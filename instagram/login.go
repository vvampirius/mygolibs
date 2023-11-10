package instagram

import (
	"errors"
	"github.com/Davincible/goinsta/v3"
)

func InstagramLogin(sessionFilePath, username, password string) (*goinsta.Instagram, error) {
	client, err := goinsta.Import(sessionFilePath)
	if err == nil {
		return client, nil
	}
	DebugLog.Println(err.Error())
	if username == `` || password == `` {
		err := errors.New(`No creds provided`)
		ErrorLog.Println(err.Error())
		return nil, err
	}
	client = goinsta.New(username, password)
	if err := client.Login(); err != nil {
		ErrorLog.Println(err.Error())
		return nil, err
	}
	return client, nil
}
