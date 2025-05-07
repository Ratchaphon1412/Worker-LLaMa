package utils

import (
	"os"
)

func SaveToFile(filePath string, data []byte) (*os.File, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return nil, err
	}

	return file, nil
}
