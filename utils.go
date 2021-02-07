package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DoesFileExist checks íf a given file exists or not
func DoesFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

// ReadFileToString returns a string with the file's content. Returns "" on error
func ReadFileToString(filePath string) string {
	var fileBytes []byte
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return string(fileBytes)

}

//GetHome returns the home directory of the current user
func GetHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err.Error())
		return ""
	}
	return home
}

// GenerateRandomString creates a random string of length n
func GenerateRandomString(n int) string {

	// Seed the time
	rand.Seed(time.Now().UnixNano())

	// Available Chars
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// Generate
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// GetLastPartOfPath returns a path's last part like a file's name
func GetLastPartOfPath(path string) string {
	return filepath.Base(path)
}

// GetNthPartOfPath returns the nth part from a path, going right to left (file to C:)
func GetNthPartOfPath(path string, n int) string {
	filePathSplit := strings.Split(path, string(filepath.Separator))
	return filePathSplit[len(filePathSplit)-n-1]
}

// ReplaceFilepathSeparator can be used to fix inconsistent separators
func ReplaceFilepathSeparator(filePath string, newSeparator string) string {
	filePath = strings.ReplaceAll(filePath, "\\\\", newSeparator)
	filePath = strings.ReplaceAll(filePath, "\\", newSeparator)
	filePath = strings.ReplaceAll(filePath, "/", newSeparator)
	if strings.HasSuffix(filePath, newSeparator) {
		filePath = filePath[:len(filePath)-1]
	}

	return filePath
}

// SanitizeOutput tries to create a correct output path across user inputs and devices
func SanitizeOutput(outputPath, fileName string) string {
	// Given path is a directory
	if !strings.Contains(GetLastPartOfPath(outputPath), ".") {
		// Create dir if needed
		os.MkdirAll(outputPath, 0640)

		// Empty string
		if len(fileName) == 0 {
			fileName = GenerateRandomString(15)
		}

		// Append path with file's name
		if strings.HasSuffix(outputPath, string(filepath.Separator)) {
			outputPath += GetLastPartOfPath(fileName)
		} else {
			outputPath += string(filepath.Separator) + GetLastPartOfPath(fileName)
		}
	} else {
		// Output is a file, create it's directory if needed
		os.MkdirAll(GetNthPartOfPath(outputPath, 1), 0640)
	}

	return outputPath
}

// ZipFilesFromPaths packs every file contained in the slice into a single zip file
func ZipFilesFromPaths(filePaths []string) (*[]byte, error) {

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	zipWriter := zip.NewWriter(buf)

	// Add some files to the archive.
	for _, file := range filePaths {
		// Write the file or all the dir's content
		err := writeFileOrDirectoryToZip(file, zipWriter)
		if err != nil {
			return nil, err
		}
	}

	// Closing and checking for errors
	err := zipWriter.Close()
	if err != nil {
		return nil, err
	}

	// Return the bytes
	byts := buf.Bytes()
	buf.Reset()

	return &byts, nil
}

func writeFileOrDirectoryToZip(filePath string, w *zip.Writer) error {
	// Open the file
	openedFile, err := os.Open(filePath)
	if err != nil {
		return errors.New("Couldn't read File: " + filePath + " | " + err.Error())
	}

	// Given "file" is a directory
	info, err := openedFile.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {

		// Read every file
		files, err := ioutil.ReadDir(filePath)
		if err != nil {
			return errors.New("Couldn't read Dir: " + filePath + " | " + err.Error())
		}
		for _, f := range files {
			innerFilePath := filePath + string(filepath.Separator) + f.Name()

			if f.IsDir() {
				// Inner File is also a Dir
				err = writeFileOrDirectoryToZip(innerFilePath, w)
				if err != nil {
					return err
				}
			} else {
				err = writeFileToZip(innerFilePath, filePath, w)
				if err != nil {
					return err
				}
			}
		}
	} else {
		// Given File is really just a file
		writeFileToZip(filePath, filePath, w)
	}

	return nil
}

func writeFileToZip(innerFilePath string, filePath string, w *zip.Writer) error {
	// Remove relative part components
	innerFilePath = ReplaceFilepathSeparator(innerFilePath, string(filepath.Separator))
	innerFilePath = strings.ReplaceAll(innerFilePath, ".."+string(filepath.Separator), "")
	// Write File to Zip
	zipFile, err := w.Create(innerFilePath)
	if err != nil {
		return errors.New("Couldnt create File in zip: " + err.Error())
	}
	zipFileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.New("Couldnt read File: " + filePath + " | " + err.Error())
	}
	zipFile.Write(zipFileContent)
	return err
}
