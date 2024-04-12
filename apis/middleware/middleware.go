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

	// marshalledPublicKey, err := x509.MarshalPKIXPublicKey(publicKey)
	// if err != nil {
	// 	log.Fatal("Cannot marshal public key", err)
	// }

	return publicKey.(*rsa.PublicKey)

}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		publicKey := loadPublicKey()

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

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
			fmt.Println(claims["user_id"], "user_id")
			fmt.Println(ctx, "ctx")
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
	})
}
