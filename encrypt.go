package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

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
func EncryptFile(filePath string, outputPath string, useArmor bool, recipients []age.Recipient) (output string, err error) {
	// Fix potential file path errors
	filePath = ReplaceFilepathSeparator(filePath, string(filepath.Separator))

	// Read the file's content
	fileData, err := ioutil.ReadFile(filePath)

	// Error check
	if err != nil {
		return "", errors.New("inputPathError")
	}

	// Encrypt and return
	return Encrypt(&fileData, GetLastPartOfPath(filePath), outputPath, useArmor, recipients)
}

// Encrypt encrypts the given data according to the supplied args
func Encrypt(data *[]byte, fileName string, outputPath string, useArmor bool, recipients []age.Recipient) (output string, err error) {

	// Prepare the file to write to
	fullPath := GetFullPath(fileName, outputPath)
	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	// Use Armor if the user wants it
	var w, a io.WriteCloser
	if useArmor {
		a = armor.NewWriter(f)
		w, err = age.Encrypt(a, recipients...)
	} else {
		w, err = age.Encrypt(f, recipients...)
	}

	// Error check
	if err != nil {
		return "", errors.New("invalidKeyError")
	}

	// Write bytes
	_, err = io.Copy(w, bytes.NewBuffer(*data))
	if err != nil {
		return "", err
	}

	// Close
	w.Close()
	if a != nil {
		a.Close()
	}
	f.Close()

	return fullPath, nil
}
