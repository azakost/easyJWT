package easyWeb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
)

// ReadBody consumes request body, fills struct's fields and returns regex
// mismatch check boolean
func ReadBody(r *http.Request, data interface{}) bool {
	body, errorRead := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err(errorRead)
	errorUnmarshal := json.Unmarshal(body, &data)
	err(errorUnmarshal)
	regexWell := true
	fieldVals := reflect.ValueOf(data).Elem()
	fieldTags := reflect.TypeOf(data).Elem()
	for i := 0; i < fieldVals.NumField(); i++ {
		value := fieldVals.Field(i).String()
		regex := fieldTags.Field(i).Tag.Get("regex")
		if !regexp.MustCompile(regex).MatchString(value) {
			regexWell = false
			fieldVals.Field(i).SetString("regex mismatch")
		}
	}
	return regexWell
}

func WriteAsJSON(w http.ResponseWriter, code int, d interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	js, jsonError := json.Marshal(d)
	err(jsonError)
	_, writeError := w.Write(js)
	err(writeError)
}

func err(e error) {
	if e != nil {
		panic(e)
	}
}
