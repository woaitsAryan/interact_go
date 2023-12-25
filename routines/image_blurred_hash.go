package routines

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	DataURL string `json:"data_url"`
}

func GetImageBlurHash(c *fiber.Ctx, fieldName string, project *models.Project) {
	form, err := c.MultipartForm()
	if err != nil {
		return
	}

	files := form.File[fieldName]
	if files == nil {
		return
	}

	file := files[0]

	// Create a buffer to store the file content
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Create a form file field
	fileWriter, err := writer.CreateFormFile("image", file.Filename)
	if err != nil {
		helpers.LogDatabaseError("Error Getting Image Hash", err, "go_routine")
		return
	}

	// Open the file and copy its content to the form file field
	src, err := file.Open()
	if err != nil {
		helpers.LogDatabaseError("Error Getting Image Hash", err, "go_routine")
		return
	}
	defer src.Close()

	_, err = io.Copy(fileWriter, src)
	if err != nil {
		helpers.LogDatabaseError("Error Getting Image Hash", err, "go_routine")
		return
	}

	// Close the multipart writer
	writer.Close()

	URL := initializers.CONFIG.ML_URL + config.IMAGE_BLUR_HASH

	// Create a POST request to the target URL
	request, err := http.NewRequest("POST", URL, &buffer)
	if err != nil {
		helpers.LogDatabaseError("Error Getting Image Hash", err, "go_routine")
		return
	}

	// Set the Content-Type header for the request
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		helpers.LogDatabaseError("Error Getting Image Hash", err, "go_routine")
		return
	}
	defer response.Body.Close()

	// Decode the response body
	var pythonResponse Response
	if err := json.NewDecoder(response.Body).Decode(&pythonResponse); err != nil {
		helpers.LogDatabaseError("Error Getting Image Hash", err, "go_routine")
		return
	}

	// Check if the status is 'success' and return the data URL
	if pythonResponse.Status == "success" {
		project.BlurHash = pythonResponse.DataURL

		result := initializers.DB.Save(&project)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Project-GetImageBlurHash", result.Error, "go_routine")
		}
	} else {
		helpers.LogDatabaseError(fmt.Sprintf("Error Getting Image Hash - Error from Python Server %s", pythonResponse.Message), nil, "go_routine")
	}
}
