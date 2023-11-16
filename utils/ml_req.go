package utils

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
)

func MLReq(id string, url string, args ...int) ([]string, error) {

	var limit, page int
	switch len(args) {
	case 0:
		limit = 10
		page = 1
	case 1:
		limit = args[0]
		page = 1
	case 2:
		limit = args[0]
		page = args[1]
	}

	URL := initializers.CONFIG.ML_URL + url
	reqBody, _ := json.Marshal(map[string]any{
		"id":    id,
		"limit": limit,
		"page":  page,
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
