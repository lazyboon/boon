package access

import (
	"bytes"
	"encoding/json"
)

type RequestEntry struct {
	Header map[string]interface{} `json:"header,omitempty"`
	Body   interface{}            `json:"body,omitempty"`
}

type ResponseEntry struct {
	Header map[string]interface{} `json:"header,omitempty"`
	Body   interface{}            `json:"body,omitempty"`
	Status int                    `json:"status,omitempty"`
}

type Entry struct {
	Method     string         `json:"method,omitempty"`
	Path       string         `json:"path,omitempty"`
	RemoteAddr string         `json:"remote_addr,omitempty"`
	Proto      string         `json:"proto,omitempty"`
	Request    *RequestEntry  `json:"request,omitempty"`
	Response   *ResponseEntry `json:"response,omitempty"`
	Latency    string         `json:"latency,omitempty"`
	RequestId  string         `json:"request_id,omitempty"`
}

func (e *Entry) UnescapeHtmlJson() (string, error) {
	buffer := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buffer)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(e)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (e *Entry) MustUnescapeHtmlJson() string {
	ans, _ := e.UnescapeHtmlJson()
	return ans
}

func (e *Entry) Json() (string, error) {
	raw, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (e *Entry) MustJson() string {
	ans, _ := e.Json()
	return ans
}
