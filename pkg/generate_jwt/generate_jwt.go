package generatejwt

import (
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
	"os"

	"github.com/golang-jwt/jwt"
)

func loadSignKey() string {
	// Ideally this key will be fetched from a secure location such as AWS Secrets Manager or a secure file storage
	// This is just an example showing how the private key can be used to sign the JWT token
	keyData, err := os.Open("jwtRS256.key")
	if err != nil {
		log.Fatal("Cannot read private key file", err)
	}
	defer keyData.Close()

	file, err := io.ReadAll(keyData)
	if err != nil {
		log.Fatal("Cannot read public key file", err)
	}

	block, _ := pem.Decode(file)
	if block == nil {
		log.Fatal("Cannot decode key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal("Cannot parse private key", err)
	}

	privateKeyString := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}))

	return privateKeyString
}

// asymmetric encryption using a private key and a public key to verify the token
// GenerateJWT generates a JWT token using the user's email and user ID
func GenerateJWT(userEmail string, userId int) string {
	var (
		key string
		t   *jwt.Token
		s   string
	)

	key = loadSignKey()

	// Create a new token object, specifying the signing method and the claims
	// we can use the claims to store any information we want to include in the token
	t = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"email":   userEmail,
		"user_id": userId,
	})

	// Parse the private key
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
	if err != nil {
		log.Fatal("Cannot parse private key", err)
	}

	// Sign the token with the private key
	s, err = t.SignedString(privateKey)
	if err != nil {
		log.Fatal("Cannot sign key", err)
	}

	// Return the signed token
	return s
}
