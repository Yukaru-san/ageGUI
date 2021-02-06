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
)

// DecryptFile takes the given file and decrypts it's content
func DecryptFile(inputPath string, outputPath string) (output string, err error) {
	// Read the input file
	f, err := os.Open(inputPath)
	if err != nil {
		return "", errors.New("inputPathError")
	}

	// Prepare reader for one or more Identities
	var r io.Reader
	if len(Identities) == 0 {
		r, err = age.Decrypt(f, Identity)
	} else {
		r, err = age.Decrypt(f, Identities...)
	}
	if err != nil {
		return "", err
		// return errors.New("invalidKeyError")
	}

	return finishDecryptionAndSafeToFile(GetLastPartOfPath(inputPath), outputPath, r)
}

// DecryptFileWithPassword takes the given file and decrypts it's content
func DecryptFileWithPassword(inputPath string, outputPath string, password string) (output string, err error) {

	// Read the input file
	f, err := os.Open(inputPath)
	if err != nil {
		return "", errors.New("inputPathError")
	}

	// Create Identity
	i, err := age.NewScryptIdentity(password)
	if err != nil {
		return "", errors.New("invalidPasswordError")
	}

	// Prepare reader for one or more Identities
	var r io.Reader
	r, err = age.Decrypt(f, i)
	if err != nil {
		return "", errors.New("invalidPasswordError")
	}

	return finishDecryptionAndSafeToFile(GetLastPartOfPath(inputPath), outputPath, r)
}

func finishDecryptionAndSafeToFile(fileName string, outputPath string, r io.Reader) (output string, err error) {
	// Read and decrypt data
	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return "", errors.New("inputPathError")
	}

	// Sanitize Output
	if len(outputPath) == 0 {
		outputPath = GetHome() + string(filepath.Separator) + "age" + string(filepath.Separator) + "decrypted"
		os.MkdirAll(outputPath, 0644)
	}
	outputPath = SanitizeOutput(outputPath, fileName)

	if strings.HasSuffix(outputPath, ".enc") {
		outputPath = outputPath[:len(outputPath)-4]
	}

	// Save as file on disk
	err = ioutil.WriteFile(outputPath, out.Bytes(), 0644)
	if err != nil {
		return "", errors.New("writeError%" + err.Error())
	}

	return outputPath, err
}
