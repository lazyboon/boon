package hashid

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/speps/go-hashids/v2"
	"sync"
)

type Config struct {
	MinLength int    `json:"min_length"`
	Salt      string `json:"salt"`
	Alphabet  string `json:"alphabet"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	hashIDOnce     sync.Once
	hashIDInstance *hashids.HashID
	config         *Config
)

func InitHashID(cfg *Config) {
	config = cfg
}

func SharedHashID() *hashids.HashID {
	hashIDOnce.Do(func() {
		d := hashids.NewData()
		if config != nil {
			if config.MinLength != 0 {
				d.MinLength = 8
			}
			if config.Salt != "" {
				d.Salt = config.Salt
			}
			if config.Alphabet != "" {
				d.Alphabet = config.Alphabet
			}
		}
		hashIDInstance, _ = hashids.NewWithData(d)
	})
	return hashIDInstance
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type HashID struct {
	String string
	Int    int64
}

func NewHashID(num int64) (*HashID, error) {
	ciphertext, err := SharedHashID().Encode([]int{int(num)})
	if err != nil {
		return nil, err
	}
	return &HashID{String: ciphertext, Int: num}, nil
}

func MustNewHashID(num int64) *HashID {
	h, _ := NewHashID(num)
	return h
}

func (h *HashID) Scan(src interface{}) error {
	t, err := NewHashID(src.(int64))
	if err != nil {
		return err
	}
	*h = *t
	return nil
}

func (h *HashID) Value() (driver.Value, error) {
	return h.Int, nil
}

func (h *HashID) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String)
}

func (h *HashID) UnmarshalJSON(bytes []byte) error {
	if bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	h.String = string(bytes)
	d, err := SharedHashID().DecodeWithError(h.String)
	if err != nil {
		return err
	}
	h.Int = int64(d[0])
	return nil
}
