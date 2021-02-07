package main

import (
	"errors"
	"os"

	"filippo.io/age"
)

// ...
var (
	Recipients []age.Recipient
	Identities []age.Identity
)

// PrepareRecipient sets the Recipient up for a public key
func PrepareRecipient(publicKey string) (pubKey string, privateKey string, err error) {
	// If key is a file reference
	if DoesFileExist(publicKey) {
		f, _ := os.Open(publicKey)
		Recipients, err = age.ParseRecipients(f)
		return
	}

	// If no key was supplied
	if publicKey == "" {
		publicKey, privateKey, err = GenerateX25519Identity()
		pubKey = publicKey
	}

	// Parse the recipient
	recipient, err := age.ParseX25519Recipient(publicKey)
	if err != nil {
		return "", "", err
	}
	Recipients = append(Recipients, recipient)

	// Return error and keys if needed
	return pubKey, privateKey, err
}

// PrepareIdentity sets the Identity up for a private key
func PrepareIdentity(privateKey string) (err error) {
	// If key is a file reference
	if DoesFileExist(privateKey) {
		return PrepareAndParseIdentities(privateKey)
	}
	var identity *age.X25519Identity
	identity, err = age.ParseX25519Identity(privateKey)
	if err != nil {
		err = errors.New("invalidKeyError")
	}
	Identities = append(Identities, identity)

	return err
}

// PrepareAndParseIdentities sets the Identity up for every private key in the given file
func PrepareAndParseIdentities(keyStoragePath string) (err error) {

	// Open File
	keyFile, err := os.Open(keyStoragePath)
	if err != nil {
		return err
	}
	// Parse identites
	Identities, err = age.ParseIdentities(keyFile)
	return err
}

// GenerateX25519Identity generates and returns a generated public / private key combination
func GenerateX25519Identity() (publicKey, privateKey string, err error) {
	Identity, err := age.GenerateX25519Identity()
	if err != nil {
		return "", "", err
	}

	return Identity.Recipient().String(), Identity.String(), err
}
