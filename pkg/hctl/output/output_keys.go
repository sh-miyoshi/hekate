package output

import (
	"encoding/json"
	"fmt"

	keysapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/keys"
)

// KeysFormat ...
type KeysFormat struct {
	keys *keysapi.KeysGetResponse
}

// NewKeysFormat ...
func NewKeysFormat(keys *keysapi.KeysGetResponse) *KeysFormat {
	return &KeysFormat{
		keys: keys,
	}
}

// ToText ...
func (f *KeysFormat) ToText() (string, error) {
	res := fmt.Sprintf("Type:       %s\n", f.keys.Type)
	res += fmt.Sprintf("Public Key: %s\n", f.keys.PublicKey)

	return res, nil
}

// ToJSON ...
func (f *KeysFormat) ToJSON() (string, error) {
	bytes, err := json.Marshal(f.keys)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
