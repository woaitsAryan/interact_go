package utils

import "github.com/gofiber/fiber/v2"

func SaveMultipleFiles(c *fiber.Ctx, fieldName string, path string, resize bool, d1 int, d2 int) ([]string, error) {
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
		filePath := "public/" + path + "/" + c.GetRespHeader("loggedInUserID") + "-" + file.Filename
		c.SaveFile(file, filePath)

		if resize {
			picName, err := ResizeImage(filePath, d1, d2)
			if err != nil {
				return []string{}, err
			}
			savedImages = append(savedImages, picName)
		} else {
			savedImages = append(savedImages, filePath)
		}
	}

	return savedImages, nil
}
