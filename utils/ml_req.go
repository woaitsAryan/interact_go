package utils

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
)

func MLRecommendationsReq(id string, url string, args ...int) ([]string, error) {
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
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&responseBody); err != nil {
		return nil, helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	return responseBody.Recommendations, nil
}

func MLFlagReq(content string) (bool, error) {
	URL := initializers.CONFIG.ML_URL + config.FLAG
	reqBody, _ := json.Marshal(map[string]any{
		"content": content,
	})
	response, err := http.Post(URL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return false, helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}
	defer response.Body.Close()

	var responseBody struct {
		Flag bool `json:"flag"`
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&responseBody); err != nil {
		return false, helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	return responseBody.Flag, nil
}
