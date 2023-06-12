package utils

import "github.com/gofiber/fiber/v2"

func SaveFile(c *fiber.Ctx, fieldName string, path string, single bool, resize bool) (string, error) {

	form, err := c.MultipartForm()
	if err != nil {
		return "", err
	}

	files := form.File[fieldName]
	if single {

		if files == nil {
			return "", nil
		}

		file := files[0]

		filePath := "public/" + path + "/" + file.Filename

		c.SaveFile(file, filePath)

		if resize {
			picName, err := ResizeImage(filePath, 500, 500)
			if err != nil {
				return "", err
			}
			return picName, nil
		}

		return filePath, nil
	}

	return "", nil

}
