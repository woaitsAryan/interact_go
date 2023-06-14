package utils

import (
	"github.com/gofiber/fiber/v2"
)

func SaveFile(c *fiber.Ctx, fieldName string, path string, resize bool, d1 int, d2 int) (string, error) {

	form, err := c.MultipartForm()
	if err != nil {
		return "", err
	}

	files := form.File[fieldName]

	if files == nil {
		return "", nil
	}

	file := files[0]

	filePath := "public/" + path + "/" + c.GetRespHeader("loggedInUserID") + "-" + file.Filename

	c.SaveFile(file, filePath)

	if resize {
		picName, err := ResizeImage(filePath, d1, d2)
		if err != nil {
			return "", err
		}
		return picName, nil
	}

	return filePath, nil

}
