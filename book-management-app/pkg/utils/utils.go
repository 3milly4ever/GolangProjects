package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//function that will parse body from JSON to Go to work with it.
//use r to access the request that we will receive from the user.
//the request will have the body of the book like the fields of it title, author, etc.

// error is performed in this function so no need to assign error variable
func ParseBody(r *http.Request, x interface{}) {

	if body, err := ioutil.ReadAll(r.Body); err == nil {
		//if no error then start unmarshalling it or converting from JSON to go
		if err := json.Unmarshal([]byte(body), x); err != nil {
			//return if there is error but if not, the above gets exectucted the json gets unmarshalled.
			return
		}

	}
}
