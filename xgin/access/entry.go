package access

import (
	"bytes"
	"encoding/json"
)

type RequestEntity struct {
	Header map[string]interface{} `json:"header,omitempty"`
	Body   interface{}            `json:"body,omitempty"`
}

type ResponseEntity struct {
	Header map[string]interface{} `json:"header,omitempty"`
	Body   interface{}            `json:"body,omitempty"`
	Status int                    `json:"status,omitempty"`
}

type Entity struct {
	Method     string          `json:"method,omitempty"`
	Path       string          `json:"path,omitempty"`
	RemoteAddr string          `json:"remote_addr,omitempty"`
	Proto      string          `json:"proto,omitempty"`
	Request    *RequestEntity  `json:"request,omitempty"`
	Response   *ResponseEntity `json:"response,omitempty"`
	Latency    string          `json:"latency,omitempty"`
	RequestID  string          `json:"request_id,omitempty"`
}

func (e *Entity) UnescapeHtmlJson() (string, error) {
	buffer := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buffer)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(e)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (e *Entity) MustUnescapeHtmlJson() string {
	ans, _ := e.UnescapeHtmlJson()
	return ans
}

func (e *Entity) Json() (string, error) {
	raw, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (e *Entity) MustJson() string {
	ans, _ := e.Json()
	return ans
}
