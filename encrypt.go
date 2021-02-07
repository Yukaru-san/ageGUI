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

	// Buffer Setup
	r.SetWorkFactor(15)
	out := &bytes.Buffer{}
	var w io.WriteCloser

	// Use Armor Writer if needed
	if useArmor {
		armorWriter := armor.NewWriter(out)
		defer func() {
			armorWriter.Close()
		}()
		w, err = age.Encrypt(armorWriter, r)
	} else {
		w, err = age.Encrypt(out, r)
	}

	// Check for errors
	if err != nil {
		return "", errors.New("invalidPasswordError")
	}

	return finishEncryptionAndSafeToFile(data, inputName, outputPath, w, out)
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

	// Use byte encryption
	return Encrypt(&fileData, GetLastPartOfPath(filePath), outputPath, useArmor)
}

// Encrypt the given bytes and saves them to the outputPath
func Encrypt(data *[]byte, inputName string, outputPath string, useArmor bool) (output string, err error) {
	// Prepare the buffer and writer
	out := &bytes.Buffer{}

	var w io.WriteCloser

	// Use Armor if the user wants it
	if useArmor {
		armorWriter := armor.NewWriter(out)
		defer func() {
			armorWriter.Close()
		}()
		w, err = age.Encrypt(armorWriter, Recipients...)
	} else {
		w, err = age.Encrypt(out, Recipients...)
	}

	// Error check
	if err != nil {
		return "", errors.New("invalidKeyError")
	}

	// Do the Encryption
	return finishEncryptionAndSafeToFile(data, inputName, outputPath, w, out)
}

func finishEncryptionAndSafeToFile(data *[]byte, fileName string, outputPath string, w io.WriteCloser, out *bytes.Buffer) (output string, err error) {
	// Write bytes
	if _, err = w.Write(*data); err != nil {
		return "", errors.New("writeError%" + err.Error())
	}

	// Close
	if err := w.Close(); err != nil {
		return "", errors.New("writeError%" + err.Error())
	}

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

	// Save as file on disk
	outputPath += ".enc"
	err = ioutil.WriteFile(outputPath, out.Bytes(), 0640)
	if err != nil {
		return "", errors.New("writeError%" + err.Error())
	}

	return outputPath, err
}
