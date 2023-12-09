package utils

import (
	"fmt"
	"time"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/gofiber/fiber/v2"
)

func UploadMultipleImages(c *fiber.Ctx, fieldName string, client *helpers.BucketClient, width int, height int) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return []string{}, err
	}

	files := form.File[fieldName]
	if files == nil {
		return []string{}, nil
	}

	var savedImages []string
	for _, file := range files {
		resizedImgBuffer, err := ResizeFormImage(file, width, height)
		if err != nil {
			continue
		}

		timestamp := time.Now().UTC().Format(time.RFC3339)
		filePath := fmt.Sprintf("%s-%s-%s", c.GetRespHeader("loggedInUserID"), file.Filename, timestamp)
		resizedPicPath := fmt.Sprintf("%s-resized.jpg", filePath)

		err = client.UploadBucketFile(resizedImgBuffer, resizedPicPath)
		if err != nil {
			continue
		}

		savedImages = append(savedImages, resizedPicPath)
	}

	return savedImages, nil
}
