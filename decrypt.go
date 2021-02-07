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

// DecryptFileWithPassword takes the given file and decrypts it's content
func DecryptFileWithPassword(inputPath string, outputPath string, password string) (output string, err error) {

	// Create Identity
	identity, err := age.NewScryptIdentity(password)
	if err != nil {
		return "", errors.New("invalidPasswordError")
	}

	// Decrypt and return
	return DecryptFile(inputPath, outputPath, []age.Identity{identity})
}

// DecryptFile takes the given file and decrypts it's content
func DecryptFile(inputPath string, outputPath string, identities []age.Identity) (output string, err error) {
	// Read the input file
	f, err := os.Open(inputPath)
	if err != nil {
		return "", errors.New("inputPathError")
	}

	// Prepare the Decryption Reader. Use Armor if needed
	var ageReader io.Reader
	if b, _ := ioutil.ReadFile(inputPath); strings.HasPrefix(string(b), "-----BEGIN AGE ENCRYPTED FILE-----") {
		armorReader := armor.NewReader(f)
		ageReader, err = age.Decrypt(armorReader, identities...)
	} else {
		// Decrypt
		ageReader, err = age.Decrypt(f, identities...)
	}

	if err != nil {
		return "", err
		// return errors.New("invalidKeyError")
	}

	// Read and decrypt data
	out := &bytes.Buffer{}
	if _, err := io.Copy(out, ageReader); err != nil {
		return "", err
	}

	// Close File
	f.Close()

	// Sanitize Output
	if len(outputPath) == 0 {
		outputPath = GetHome() + string(filepath.Separator) + "age" + string(filepath.Separator) + "decrypted"
		os.MkdirAll(outputPath, 0750)
	}
	outputPath = SanitizeOutput(outputPath, GetLastPartOfPath(inputPath))

	// Remove .enc Suffix if needed
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
