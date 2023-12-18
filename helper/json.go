package helper

import (
	"encoding/json"
	"net/http"
)

// ReadFromRequestBody reads the JSON data from the provided HTTP
// request's body and decode it into the specified result interface{}.
//
// The function uses the json.NewDecoder from the encoding/json package
// to decode the JSON data.
func ReadFromRequestBody(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(result); err != nil {
		panic(err)
	}
}

// WriteToResponseBody writes the provided response data to the HTTP
// response writer in JSON format. It sets the "Content-Type" header
// to "application/json".
//
// The function uses the json.NewEncoder from the encoding/json package
// to encode the response data.
func WriteToResponseBody(writer http.ResponseWriter, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(response); err != nil {
		panic(err)
	}
}
