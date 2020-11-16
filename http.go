package easyWeb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ReadBody(r *http.Request, data interface{}) error {
	body, errorRead := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if errorRead != nil {
		panic(errorRead)
	}
	errorUnmarshal := json.Unmarshal(body, &data)
	if errorUnmarshal != nil {
		return errorUnmarshal
	}

	// Validate with regex
	// 1. Get regex function

	return nil
}
