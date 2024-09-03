package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
)

// load private key from /etc/secrets/crane/crane_private_key.pem
func LoadPrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := loadPrivateKeyFromFile("/etc/secrets/crane/crane_private_key.pem")
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func loadPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	// Read the private key file
	privateKeyFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Decode the PEM encoded private key
	block, _ := pem.Decode(privateKeyFile)
	if block == nil {
		return nil, err
	}

	// Parse the private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func DecryptData(ciphertext string, privateKey *rsa.PrivateKey) (string, error) {
	// Decode the base64 encoded ciphertext
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Decrypt the ciphertext
	hash := sha256.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, cipherBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
