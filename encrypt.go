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
func EncryptWithPassword(data *[]byte, inputName string, outputPath string, useArmor bool, password string) (output string, err error) {
	// Prepare PW Recipient
	r, err := age.NewScryptRecipient(password)
	if err != nil {
		return "", errors.New("invalidPasswordError")
	}

	// Set WorkFactor
	r.SetWorkFactor(15)

	// Encrypt and return
	return Encrypt(data, inputName, outputPath, useArmor, []age.Recipient{r})
}

// EncryptFile the given file and saves them to the outputPath
func EncryptFile(filePath string, outputPath string, useArmor bool) (output string, err error) {
	// Fix potential file path errors
	filePath = ReplaceFilepathSeparator(filePath, string(filepath.Separator))

	// Read the file's content
	fileData, err := ioutil.ReadFile(filePath)

	// Error check
	if err != nil {
		return "", errors.New("inputPathError")
	}

	// Encrypt and return
	return Encrypt(&fileData, GetLastPartOfPath(filePath), outputPath, useArmor, Recipients)
}

// Encrypt encrypts the given data according to the supplied args
func Encrypt(data *[]byte, fileName string, outputPath string, useArmor bool, recipients []age.Recipient) (output string, err error) {
	// Sanitize Output
	if len(outputPath) == 0 {
		outputPath = filepath.Join(GetHome(), "age", "encrypted")
		os.MkdirAll(outputPath, 0750)
	} else if !strings.Contains(outputPath, string(filepath.Separator)) {
		fileName = outputPath
		outputPath = filepath.Join(GetHome(), "age", "encrypted")
		os.MkdirAll(outputPath, 0750)
	}
	outputPath = SanitizeOutput(outputPath, fileName) + ".enc"

	// Create file on disk
	f, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}

	// Create Encryption Writer and use armor if needed
	var ageWriter, armorWriter io.WriteCloser

	if useArmor {
		armorWriter := armor.NewWriter(f)
		ageWriter, err = age.Encrypt(armorWriter, recipients...)
	} else {
		ageWriter, err = age.Encrypt(f, recipients...)
	}

	// Check for errors
	if err != nil {
		return "", err
	}

	// Write bytes
	if _, err = io.Copy(ageWriter, bytes.NewBuffer(*data)); err != nil {
		return "", err
	}

	// Close
	ageWriter.Close()
	if armorWriter != nil {
		armorWriter.Close()
	}
	f.Close()

	// Return
	return outputPath, err
}
