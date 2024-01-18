package utils

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/gofiber/fiber/v2"
)

func UploadImage(c *fiber.Ctx, fieldName string, client *helpers.BucketClient, width int, height int) (string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return "", err
	}

	files := form.File[fieldName]
	if files == nil {
		return "", nil
	}

	file := files[0]

	resizedImgBuffer, err := ResizeFormImage(file, width, height)
	if err != nil {
		return "", err
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	filePath := fmt.Sprintf("%s-%s-%s", c.GetRespHeader("loggedInUserID"), timestamp, file.Filename)
	resizedPicPath := fmt.Sprintf("%s-resized.jpg", filePath)

	err = client.UploadBucketFile(resizedImgBuffer, resizedPicPath)
	if err != nil {
		return "", err
	}

	return resizedPicPath, nil
}

func UploadResume(c *fiber.Ctx) (string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return "", err
	}

	files := form.File["resume"]
	if files == nil {
		return "", nil
	}

	file := files[0]

	fileContent, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileContent.Close()

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, fileContent); err != nil {
		return "", err
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	filePath := fmt.Sprintf("%s-%s-%s", c.GetRespHeader("loggedInUserID"), timestamp, file.Filename)

	err = helpers.UserResumeClient.UploadBucketFile(&buffer, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func UploadFile(c *fiber.Ctx) (string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return "", err
	}

	files := form.File["file"]
	if files == nil {
		return "", fmt.Errorf("file not present")
	}

	file := files[0]

	fileContent, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileContent.Close()

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, fileContent); err != nil {
		return "", err
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	filePath := fmt.Sprintf("%s-%s-%s", c.Params("orgID"), timestamp, SoftSlugify(file.Filename))

	err = helpers.ResourceClient.UploadBucketFile(&buffer, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
