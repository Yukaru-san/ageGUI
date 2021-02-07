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

// DecryptFile takes the given file and decrypts it's content
func DecryptFile(inputPath string, outputPath string) (output string, err error) {
	// Read the input file
	f, err := os.Open(inputPath)
	if err != nil {
		return "", errors.New("inputPathError")
	}

	// Prepare reader for one or more Identities
	var r io.Reader

	// When the file is using Armor TODO make it more efficient TODO Make it work at all lol
	if b, _ := ioutil.ReadFile(inputPath); strings.HasPrefix(string(b), "-----BEGIN AGE ENCRYPTED FILE-----") {
		armorReader := armor.NewReader(f)
		r, err = age.Decrypt(armorReader, Identities...)
	} else {
		// Decrypt
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

	// When the file is using Armor TODO make it more efficient TODO Make it work at all lol
	if b, _ := ioutil.ReadFile(inputPath); strings.HasPrefix(string(b), "-----BEGIN AGE ENCRYPTED FILE-----") {
		armorReader := armor.NewReader(f)
		r, err = age.Decrypt(armorReader, i)
	} else {
		// Decrypt
		r, err = age.Decrypt(f, i)
	}

	if err != nil {
		return "", errors.New("invalidPasswordError")
	}

	return finishDecryptionAndSafeToFile(GetLastPartOfPath(inputPath), outputPath, r)
}

func finishDecryptionAndSafeToFile(fileName string, outputPath string, r io.Reader) (output string, err error) {
	// Read and decrypt data
	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return "", errors.New("io.Copy Error: " + err.Error())
		// return "", errors.New("inputPathError")
	}

	// Sanitize Output
	if len(outputPath) == 0 {
		outputPath = GetHome() + string(filepath.Separator) + "age" + string(filepath.Separator) + "decrypted"
		os.MkdirAll(outputPath, 0750)
	}
	outputPath = SanitizeOutput(outputPath, fileName)

	if strings.HasSuffix(outputPath, ".enc") {
		outputPath = outputPath[:len(outputPath)-4]
	}

	// Save as file on disk
	err = ioutil.WriteFile(outputPath, out.Bytes(), 0640)
	if err != nil {
		return "", errors.New("writeError%" + err.Error())
	}

	return outputPath, err
}
