package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
	"filippo.io/age/armor"
)

// EncryptFileWithPassword the given file and saves them to the outputPath
func EncryptFileWithPassword(filePath string, outputPath string, useArmor bool, password string) (output string, err error) {
	// Read the file's content
	fileData, err := ioutil.ReadFile(filePath)

	// Error check
	if err != nil {
		return "", errors.New("inputPathError")
	}

	return EncryptWithPassword(&fileData, GetLastPartOfPath(filePath), outputPath, useArmor, password)
}

// EncryptWithPassword the given bytes and saves them to the outputPath
func EncryptWithPassword(data *[]byte, fileName string, outputPath string, useArmor bool, password string) (output string, err error) {
	fullPath := getFullPath(fileName, outputPath)
	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	// Prepare PW Recipient
	r, err := age.NewScryptRecipient(password)
	if err != nil {
		return "", errors.New("invalidPasswordError")
	}

	// Buffer Setup
	r.SetWorkFactor(15)
	var w, a io.WriteCloser

	// Use Armor Writer if needed
	if useArmor {
		a = armor.NewWriter(f)
		w, err = age.Encrypt(a, r)
	} else {
		w, err = age.Encrypt(f, r)
	}

	// Check for errors
	if err != nil {
		return "", errors.New("invalidPasswordError")
	}

	_, err = io.Copy(w, bytes.NewBuffer(*data))
	if err != nil {
		return "", err
	}

	w.Close()
	if a != nil {
		a.Close()
	}
	f.Close()

	return fullPath, nil
}

// EncryptFile the given file and saves them to the outputPath
func EncryptFile(filePath string, outputPath string, useArmor bool, Recipients []age.Recipient) (output string, err error) {
	// Fix potential file path errors
	filePath = ReplaceFilepathSeparator(filePath, string(filepath.Separator))

	// Read the file's content
	fileData, err := ioutil.ReadFile(filePath)

	// Error check
	if err != nil {
		return "", errors.New("inputPathError")
	}

	// Use byte encryption
	return Encrypt(&fileData, GetLastPartOfPath(filePath), outputPath, useArmor, Recipients)
}

// Encrypt the given bytes and saves them to the outputPath
func Encrypt(data *[]byte, fileName string, outputPath string, useArmor bool, Recipients []age.Recipient) (output string, err error) {
	fullPath := getFullPath(fileName, outputPath)
	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	var w, a io.WriteCloser

	// Use Armor if the user wants it
	if useArmor {
		a = armor.NewWriter(f)
		w, err = age.Encrypt(a, Recipients...)
	} else {
		w, err = age.Encrypt(f, Recipients...)
	}
	// Error check
	if err != nil {
		return "", errors.New("invalidKeyError")
	}

	_, err = io.Copy(w, bytes.NewBuffer(*data))
	if err != nil {
		return "", err
	}

	w.Close()
	if a != nil {
		a.Close()
	}
	f.Close()
	return outputPath, nil
}

func getFullPath(fileName string, outputPath string) string {
	// Sanitize Output
	if len(outputPath) == 0 {
		outputPath = GetHome() + string(filepath.Separator) + "age" + string(filepath.Separator) + "encrypted"
		os.MkdirAll(outputPath, 0750)
	} else if !strings.Contains(outputPath, string(filepath.Separator)) {
		fileName = outputPath
		outputPath = GetHome() + string(filepath.Separator) + "age" + string(filepath.Separator) + "encrypted"
		os.MkdirAll(outputPath, 0750)
	}
	outputPath = SanitizeOutput(outputPath, fileName)
	outputPath += ".enc"

	return outputPath
}
