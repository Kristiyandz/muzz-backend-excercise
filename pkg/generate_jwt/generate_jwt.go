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
	// keyData, err := ioutil.ReadFile("jwtRS256.key")
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
func GenerateJWT(userEmail string, userId int) string {
	var (
		key string
		t   *jwt.Token
		s   string
	)

	key = loadSignKey()

	t = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"email":   userEmail,
		"user_id": userId,
	})
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
	if err != nil {
		log.Fatal("Cannot parse private key", err)
	}

	s, err = t.SignedString(privateKey)
	if err != nil {
		log.Fatal("Cannot sign key", err)
	}

	return s
}
