package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryption(t *testing.T) {

	keyName, pubName, err := generateKeys()
	require.NoError(t, err)
	defer os.Remove(keyName)
	defer os.Remove(pubName)

	key, err := PrivateKeyFromFile(keyName)
	require.NoError(t, err)

	pub, err := PublicKeyFromFile(pubName)
	require.NoError(t, err)

	data := []byte("very secret data")
	encrypted, err := Encrypt(pub, data)
	require.NoError(t, err)

	decrypted, err := Decrypt(key, encrypted)
	require.NoError(t, err)

	assert.Equal(t, data, decrypted)

}

func generateKeys() (string, string, error) {

	filename := "key"
	bitSize := 4096

	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return "", "", err
	}
	pub := key.Public()

	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)

	keyFileName := fmt.Sprintf("%s.rsa", filename)
	if err := os.WriteFile(keyFileName, keyPEM, 0700); err != nil {
		return "", "", err
	}

	pubFileName := fmt.Sprintf("%s.rsa.pub", filename)
	if err := os.WriteFile(pubFileName, pubPEM, 0755); err != nil {
		return "", "", err
	}

	return keyFileName, pubFileName, nil
}
