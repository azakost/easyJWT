package easyWeb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"time"
)

// ReadBody consumes request body, fills struct's fields and returns regex
// mismatch check boolean
func ReadBody(r *http.Request, data interface{}) bool {
	body, errorRead := ioutil.ReadAll(r.Body)
	err(errorRead)
	defer r.Body.Close()
	err(json.Unmarshal(body, &data))
	regexWell := true
	fieldVals := reflect.ValueOf(data).Elem()
	fieldTags := reflect.TypeOf(data).Elem()
	for i := 0; i < fieldVals.NumField(); i++ {
		value := fieldVals.Field(i).String()
		regex := fieldTags.Field(i).Tag.Get("regex")
		if !regexp.MustCompile(regex).MatchString(value) {
			regexWell = false
			fieldVals.Field(i).SetString("regex!")
		}
	}
	return regexWell
}

func ReadParams(r *http.Request, data interface{}) bool {
	regexWell := true
	fieldVals := reflect.ValueOf(data).Elem()
	fieldTags := reflect.TypeOf(data).Elem()
	for i := 0; i < fieldVals.NumField(); i++ {
		name := fieldTags.Field(i).Tag.Get("json")
		if value, ok := r.URL.Query()[name]; ok {
			regex := fieldTags.Field(i).Tag.Get("regex")
			if !regexp.MustCompile(regex).MatchString(value[0]) {
				regexWell = false
				fieldVals.Field(i).SetString("regex!")
			} else {
				fieldVals.Field(i).SetString(value[0])
			}
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

func PutCookie(w http.ResponseWriter, name, value string, exp time.Time) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  exp,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}
