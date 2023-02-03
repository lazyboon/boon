package cid

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

var cipherKey string

func InjectKey(key string) {
	cipherKey = key
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type CID struct {
	// encrypted string
	// the algorithm used is cipher.NewCBCEncrypter, and the result is base64.RawURLEncoding
	Ciphertext string

	// id of the original numeric type
	Original int64
}

func NewWithOriginal(original int64) (*CID, error) {
	ciphertext, err := encrypt([]byte(fmt.Sprintf("%d", original)), []byte(cipherKey))
	if err != nil {
		return nil, err
	}
	return &CID{Ciphertext: ciphertext, Original: original}, nil
}

func NewWithCiphertext(ciphertext string, key string) (*CID, error) {
	bts, err := decrypt([]byte(ciphertext), []byte(key))
	if err != nil {
		return nil, err
	}
	original, err := strconv.ParseInt(string(bts), 10, 64)
	if err != nil {
		return nil, err
	}
	return &CID{Ciphertext: ciphertext, Original: original}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Scan implement sql.Scanner
func (c *CID) Scan(src interface{}) error {
	t, err := NewWithOriginal(src.(int64))
	if err != nil {
		return err
	}
	*c = *t
	return nil
}

// Value implement driver.Valuer
func (c CID) Value() (driver.Value, error) {
	return c.Original, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// MarshalJSON implement json.Marshaler
func (c *CID) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Ciphertext)
}

// UnmarshalJSON implement json.Unmarshaler
func (c *CID) UnmarshalJSON(bytes []byte) error {
	if bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	c.Ciphertext = string(bytes)

	d, err := decrypt(bytes, []byte(cipherKey))
	if err != nil {
		return err
	}
	i, err := strconv.ParseInt(string(d), 10, 64)
	if err != nil {
		return err
	}
	c.Original = i

	return nil
}
