package utils

import (
	"fmt"
	"time"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/gofiber/fiber/v2"
)

func UploadFile(c *fiber.Ctx, fieldName string, client *helpers.BucketClient, d1 int, d2 int) (string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return "", err
	}

	files := form.File[fieldName]
	if files == nil {
		return "", nil
	}

	file := files[0]

	resizedImgBuffer, err := ResizeFormImage(file, d1, d2)
	if err != nil {
		return "", err
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	filePath := fmt.Sprintf("%s-%s-%s", c.GetRespHeader("loggedInUserID"), file.Filename, timestamp)
	resizedPicPath := fmt.Sprintf("%s-resized.jpg", filePath)

	err = client.UploadBucketFile(resizedImgBuffer, resizedPicPath)
	if err != nil {
		return "", err
	}

	return resizedPicPath, nil
}