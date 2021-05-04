package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Encode_to_json(w http.ResponseWriter, body interface{}) {
	w.Header().Add("Content-Type", "application/json")
	jsonDataByte, _ := json.Marshal(body)
	jsonData := string(jsonDataByte)
	_, _ = fmt.Fprint(w, jsonData)
}
