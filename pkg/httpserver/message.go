package httpserver

import (
	"encoding/json"
)

type jsonMessage struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func Message(status string, data interface{}) string {
	resp := &jsonMessage{status, data}
	out, _ := json.Marshal(resp)
	return string(out)
}
