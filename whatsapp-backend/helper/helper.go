package helper

import (
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GenerateUniqueFilename(originalFilename string) string {
	// Extract the file extension
	fileExt := filepath.Ext(originalFilename)

	// Generate a random string
	rand.Seed(time.Now().UnixNano())
	randomString := strconv.FormatInt(rand.Int63(), 10)

	// Create a unique filename with timestamp and random string
	filename := time.Now().Format("20060102150405") + "_" + randomString + fileExt

	// Remove any special characters or spaces
	filename = strings.ReplaceAll(filename, " ", "_")

	return filename

}
