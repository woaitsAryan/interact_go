package utils

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
)

func MLReq(id string, url string) ([]string, error) {
	URL := "http://127.0.0.1:3030/" + url
	reqBody, _ := json.Marshal(map[string]string{
		"id": id,
	})
	response, err := http.Post(URL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}
	defer response.Body.Close()

	var responseBody struct {
		Recommendations []string `json:"recommendations"`
		// Define other fields you expect in the response
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&responseBody); err != nil {
		return nil, helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	return responseBody.Recommendations, nil
}
