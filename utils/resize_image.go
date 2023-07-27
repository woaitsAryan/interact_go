package utils

import (
	"image/jpeg"
	"os"
	"path"
	"strings"

	"github.com/nfnt/resize"
)

func ResizeImage(picPath string, d1 int, d2 int) (string, error) {
	file, err := os.Open(picPath)
	if err != nil {
		return "", err
	}

	img, err := jpeg.Decode(file) //! Does not works with jpg format
	if err != nil {
		return "", err
	}
	file.Close()

	m := resize.Resize(uint(d1), uint(d2), img, resize.Lanczos3)

	extension := path.Ext(picPath)
	fileNameWithoutExt := picPath[:len(picPath)-len(extension)]

	resizedPicPath := fileNameWithoutExt + "-resized.jpg" //!add time in the name

	out, err := os.Create(resizedPicPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	err = jpeg.Encode(out, m, nil)
	if err != nil {
		return "", err
	}

	err = os.Remove(picPath)
	if err != nil {
		return "", err
	}

	resizedPicArr := strings.Split(resizedPicPath, "/")

	resizedPic := resizedPicArr[len(resizedPicArr)-1]

	return resizedPic, nil
}
