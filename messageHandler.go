package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
)

// HandleMessage does stuff with a msg from js returns a response
func HandleMessage(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "getBaseDirectory":
		// Send the current working directory to JS
		dir, _ := os.Getwd()
		payload = dir
		return
	case "ageRequest":
		// Parse the request
		var request AgeRequest
		err = json.Unmarshal(m.Payload, &request)

		// Return on error
		if err != nil {
			payload = "JSON Error"
			return
		}

		// Handle the request and answer if needed
		var publicKey, privateKey, outputPath string
		publicKey, privateKey, outputPath, err = handleAgeRequest(request)

		if err != nil {
			Logger.Println("\n\nERROR OCCURED: " + err.Error() + "\n\n")
			payload = err.Error()
			return
		} else if len(publicKey) > 0 {
			payload = "generatedKeys%" + publicKey + "%" + privateKey + "%" + outputPath
		} else if len(outputPath) > 0 {
			payload = "outputPath%" + outputPath
		}

		return
	default:
		// Unknown command
		payload = "Unknown command"
		return
	}

	return
}

// Parses a given AgeRequest, does what it needs to do, and returns success / error
func handleAgeRequest(request AgeRequest) (publicKey string, privateKey string, outputPath string, err error) {
	// Fix potential file path errors in output
	request.OutputPath = ReplaceFilepathSeparator(request.OutputPath, string(filepath.Separator))

	// Age Key Encryption / Decryption Setup
	if !request.UsePassword && request.Encrypt {
		publicKey, privateKey, err = PrepareRecipient(request.CryptKey)
		if err != nil {
			return "", "", "", err
		}
	} else if !request.UsePassword {
		err = PrepareIdentity(request.CryptKey)
		if err != nil {
			return "", "", "", err
		}
	}

	// Encrypt and Zip Files if needed
	if request.Encrypt && request.ZipFiles {
		// Prepare Zip
		zipBytes, err := ZipFilesFromPaths(request.Files)
		if err != nil {
			return "", "", "", err
		}

		// Password
		if request.UsePassword {
			outputPath, err = EncryptWithPassword(zipBytes, GenerateRandomString(10)+".zip", request.OutputPath, request.UseArmor, request.CryptKey)
			if err != nil {
				return "", "", "", err
			}
		} else {
			// Age Key
			outputPath, err = Encrypt(zipBytes, GenerateRandomString(10)+".zip", request.OutputPath, request.UseArmor, Recipients)
			if err != nil {
				return "", "", "", err
			}
		}
	} else {
		// Iterate all File paths
		for _, file := range request.Files {
			// Fix potential file path errors
			file = ReplaceFilepathSeparator(file, string(filepath.Separator))

			// Password being used
			if request.UsePassword {
				if request.Encrypt {
					outputPath, err = EncryptFileWithPassword(file, request.OutputPath, request.UseArmor, request.CryptKey)
				} else {
					outputPath, err = DecryptFileWithPassword(file, request.OutputPath, request.CryptKey)
				}
			} else {
				// Age Key being used
				if request.Encrypt {
					outputPath, err = EncryptFile(file, request.OutputPath, request.UseArmor)
				} else {
					outputPath, err = DecryptFile(file, request.OutputPath, Identities)
				}
			}
		}
	}

	// JS Seite: Datei per Auswahl hinzufügen können, oder Drag & Drop
	if publicKey == request.CryptKey {
		publicKey = ""
		privateKey = ""
	}

	return publicKey, privateKey, outputPath, err
}

// AgeRequest contains all info needed to proceed
type AgeRequest struct {
	Encrypt     bool     `json:"encrypt"`
	ZipFiles    bool     `json:"zip"`
	UseArmor    bool     `json:"armor"`
	CryptKey    string   `json:"key"`
	UsePassword bool     `json:"usePassword"`
	OutputPath  string   `json:"output"`
	Files       []string `json:"paths"`
}
