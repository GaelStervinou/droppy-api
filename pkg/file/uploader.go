package file

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func UploadFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("unable to open file: %v", err)
	}
	defer src.Close()

	// Create the target directory if it doesn't exist
	targetDir := "assets"
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("unable to create directory: %v", err)
	}

	// Generate a unique filename
	uniqueID := uuid.New().String()
	extension := filepath.Ext(file.Filename)
	uniqueFilename := uniqueID + extension
	targetFilePath := filepath.Join(targetDir, uniqueFilename)

	// Create the target file
	dst, err := os.Create(targetFilePath)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %v", err)
	}
	defer dst.Close()

	// Copy the file content to the target file
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("unable to save file: %v", err)
	}

	return os.Getenv("BASE_URL") + "/" + targetFilePath, nil
}
