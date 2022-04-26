package cora

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestResponse_JSON(t *testing.T) {
	var jsonMap map[string]interface{}
	// json input data
	in := "{\"foo\":{\"baz\": [1,2,3]}}"
	// create an array of response to test
	responseTests := []Response{
		{200, nil, []byte(in), &http.Response{Status: "200 OK"}},
		{400, errors.New("error occurred"), []byte{}, &http.Response{StatusCode: 400}},
	}
	// convert to json
	for _, elm := range responseTests {
		err := json.Unmarshal(elm.Body, &jsonMap)
		if err != nil && elm.Err != nil {
			if elm.Code != elm.Raw.StatusCode || len(elm.Body) > 0 {
				t.Fatal(err)
			}
		} else {
			var expected map[string]interface{}
			err = json.Unmarshal(elm.Body, &expected)
			if err != nil {
				t.Log(err)
			}
			if !reflect.DeepEqual(jsonMap, expected) {
				t.Log(err)
			}
		}
	}
}
