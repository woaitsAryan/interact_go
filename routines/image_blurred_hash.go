package routines

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"

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

func processCtx(c *fiber.Ctx, fieldName string) ([]*multipart.FileHeader, error) {
	if c == nil {
		return nil, fmt.Errorf("ctx is nil")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	files := form.File[fieldName]
	return files, nil
}

func makeBlurHashReq(file *multipart.FileHeader) (*Response, error) {
	// Create a buffer to store the file content
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Create a form file field
	fileWriter, err := writer.CreateFormFile("image", file.Filename)
	if err != nil {
		return nil, err
	}

	// Open the file and copy its content to the form file field
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	_, err = io.Copy(fileWriter, src)
	if err != nil {
		return nil, err
	}

	// Close the multipart writer
	writer.Close()

	URL := initializers.CONFIG.ML_URL + config.IMAGE_BLUR_HASH

	// Create a POST request to the ml URL
	request, err := http.NewRequest("POST", URL, &buffer)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var pythonResponse Response
	if err := json.NewDecoder(response.Body).Decode(&pythonResponse); err != nil {
		return nil, err
	}

	return &pythonResponse, nil
}

func GetImageBlurHash(c *fiber.Ctx, fieldName string, model interface{}) {
	files, err := processCtx(c, fieldName)
	if err != nil {
		go helpers.LogDatabaseError("Error Getting Image Hash", err, "go_routine")
		return
	}
	if files == nil {
		return
	}

	file := files[0]

	pythonResponse, err := makeBlurHashReq(file)
	if err != nil {
		go helpers.LogDatabaseError("Error Getting Image Hash", err, "go_routine")
		return
	}

	if pythonResponse.Status == "success" {
		modelValue := reflect.ValueOf(model).Elem()
		blurHashField := modelValue.FieldByName("BlurHash")

		if blurHashField.IsValid() && blurHashField.CanSet() {
			blurHashField.SetString(pythonResponse.DataURL)

			result := initializers.DB.Save(model)
			if result.Error != nil {
				go helpers.LogDatabaseError(fmt.Sprintf("Error while updating model - GetImageBlurHash: %v", result.Error), nil, "go_routine")
			}
		} else {
			go helpers.LogDatabaseError("Invalid or unexported field", nil, "go_routine")
		}
	} else {
		go helpers.LogDatabaseError(fmt.Sprintf("Error Getting Image Hash - Error from Python Server %s", pythonResponse.Message), nil, "go_routine")
	}
}

func GetBlurHashesForPost(c *fiber.Ctx, fieldName string, post *models.Post) {
	files, err := processCtx(c, fieldName)
	if err != nil {
		go helpers.LogDatabaseError("Error Getting Image Hash", err, "go_routine")
		return
	}
	if files == nil {
		return
	}

	var hashes []string
	for _, file := range files {
		pythonResponse, err := makeBlurHashReq(file)
		if err != nil {
			go helpers.LogDatabaseError("Error Getting Image Hash", err, "go_routine")
			return
		}

		if pythonResponse.Status == "success" {
			hashes = append(hashes, pythonResponse.DataURL)
		} else {
			go helpers.LogDatabaseError(fmt.Sprintf("Error Getting Image Hash - Error from Python Server %s", pythonResponse.Message), nil, "go_routine")
		}
	}

	post.Hashes = hashes

	result := initializers.DB.Save(post)
	if result.Error != nil {
		go helpers.LogDatabaseError(fmt.Sprintf("Error while updating model - GetBlurHashesForPost: %v", result.Error), nil, "go_routine")
	}
}
