package utils

import (
	"log"
	"os"
)

func DeleteFile(path string, fileName string) error {

	if fileName == "" || fileName == "default.jpg" {
		return nil
	}

	picPath := "public/" + path + "/" + fileName

	log.Println(picPath)

	err := os.Remove(picPath)
	if err != nil {
		return err
	}

	return nil
}
