package output

import (
	"encoding/json"
	"fmt"

	clientapi "github.com/sh-miyoshi/jwt-server/pkg/apihandler/v1/client"
)

// ClientInfoFormat ...
type ClientInfoFormat struct {
	client *clientapi.ClientGetResponse
}

// ClientsInfoFormat ...
type ClientsInfoFormat struct {
	clients []*clientapi.ClientGetResponse
}

// NewClientInfoFormat ...
func NewClientInfoFormat(client *clientapi.ClientGetResponse) *ClientInfoFormat {
	return &ClientInfoFormat{
		client: client,
	}
}

// NewClientsInfoFormat ...
func NewClientsInfoFormat(clients []*clientapi.ClientGetResponse) *ClientsInfoFormat {
	return &ClientsInfoFormat{
		clients: clients,
	}
}

// ToText ...
func (f *ClientInfoFormat) ToText() (string, error) {
	res := fmt.Sprintf("ID:                  %s\n", f.client.ID)
	res += fmt.Sprintf("Secret:              %s\n", f.client.Secret)
	res += fmt.Sprintf("AccessType:          %s\n", f.client.AccessType)
	res += fmt.Sprintf("CreatedAt:           %s\n", f.client.CreatedAt)
	res += fmt.Sprintf("AllowedCallbackURLs: %v", f.client.AllowedCallbackURLs)
	return res, nil
}

// ToJSON ...
func (f *ClientInfoFormat) ToJSON() (string, error) {
	bytes, err := json.Marshal(f.client)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ToText ...
func (f *ClientsInfoFormat) ToText() (string, error) {
	res := ""
	for i, prj := range f.clients {
		format := NewClientInfoFormat(prj)
		msg, err := format.ToText()
		if err != nil {
			return "", err
		}
		res += msg
		if i < len(f.clients)-1 {
			res += "\n---\n"
		}
	}
	return res, nil
}

// ToJSON ...
func (f *ClientsInfoFormat) ToJSON() (string, error) {
	bytes, err := json.Marshal(f.clients)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
