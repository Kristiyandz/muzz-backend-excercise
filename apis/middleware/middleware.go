package middleware

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func loadPublicKey() *rsa.PublicKey {
	// Ideally his key will be fetched from a secure location such as AWS Secrets Manager or a secure file storage
	// This is just an example showing how the public key can be used to verify the JWT token
	rawFile, err := os.Open("jwtRS256.key.pub")
	if err != nil {
		log.Fatal("Cannot read public key file", err)
	}
	defer rawFile.Close()

	// Step 2: Read the file
	file, err := io.ReadAll(rawFile)
	if err != nil {
		log.Fatal("Cannot read public key file", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(file)
	if block == nil {
		log.Fatal("failed to parse PEM block containing the public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal("Cannot parse public key", err)
	}

	return publicKey.(*rsa.PublicKey)

}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Check the format of the Authorization header
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Load the public key
		publicKey := loadPublicKey()

		// Parse the JWT token
		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check if the token is valid and extract the claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])

			// Token is valid, proceed with the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			// Token is invalid
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
	})
}
