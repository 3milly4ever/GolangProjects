package utils

import (
	"encoding/json"
	"net/http"
)

func ParseBody(r *http.Request, x interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(x); err != nil {
		return err
	}
	return nil
}

//parses from JSON to Go struct

// func ParseBody(r *http.Request, x interface{}) {

// 	if body, err := ioutil.ReadAll(r.Body); err == nil {
// 		if err := json.Unmarshal([]byte(body), x); err != nil {
// 			return
// 		}
// 	}

// }

// func ParseBody(r *http.Request, x interface{}) error {
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		return err
// 	}

// 	if err := json.Unmarshal(body, x); err != nil {
// 		return err
// 	}

// 	return nil
// }
