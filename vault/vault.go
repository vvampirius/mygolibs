package vault

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/vault/api"
	"net/http"
	"time"
)

type Vault struct {
	Addr string
	Token string
	Name string
	Value string
	onChange func(old, new string)
	onError func(err error)
}

func (vault *Vault) LoadValue() error {
	client, err := api.NewClient(&api.Config{Address: vault.Addr, HttpClient: &http.Client{Timeout: 10 * time.Second}})
	if err != nil {
		ErrorLog.Println(err.Error())
		return err
	}
	client.SetToken(vault.Token)
	response, err := client.Logical().Read(vault.Name)
	if err != nil {
		ErrorLog.Println(err.Error())
		return err
	}
	if response == nil {
		err := errors.New(fmt.Sprintf("%s not found", vault.Name))
		ErrorLog.Println(err.Error())
		return err
	}
	newValue := fmt.Sprintf("%v", response.Data[`value`])
	if vault.Value == newValue { return nil }
	if vault.onChange != nil {
		vault.onChange(vault.Value, newValue)
	}
	vault.Value = newValue
	return nil
}

func (vault *Vault) ReloadRoutine(interval time.Duration, ctx context.Context) {
	for {
		time.Sleep(interval)
		select {
		case <-ctx.Done():
			return
		default:
			vault.LoadValue()
		}
	}
}


func NewVault(addr string, token string, name string) (*Vault, error) {
	vault := Vault{
		Addr: addr,
		Token: token,
		Name: name,
	}
	if err := vault.LoadValue(); err != nil { return nil, err }
	return &vault, nil
}
