package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"log"

	"github.com/AKSHAY-1505/auth-service/initializers"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func InitJWTKeys() {
	var err error

	privateKeyBase64 := initializers.AppConfig.JWTPrivateKey
	if privateKeyBase64 == "" {
		panic("private key not configured in environment variables")
	}

	publicKeyBase64 := initializers.AppConfig.JWTPublicKey
	if publicKeyBase64 == "" {
		panic("public key not configured in environment variables")
	}

	privateKey, err = parseRSAPrivateKeyFromEnv(privateKeyBase64)
	if err != nil {
		log.Fatalf("failed to parse JWT private key: %v", err)
	}

	publicKey, err = parseRSAPublicKeyFromEnv(publicKeyBase64)
	if err != nil {
		log.Fatalf("failed to parse JWT public key: %v", err)
	}
}

func GetPrivateKey() *rsa.PrivateKey {
	return privateKey
}

func GetPublicKey() *rsa.PublicKey {
	return publicKey
}

func parseRSAPrivateKeyFromEnv(b64 string) (*rsa.PrivateKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("invalid PEM block")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("key is not RSA")
	}

	return rsaKey, nil
}

func parseRSAPublicKeyFromEnv(publicKeyB64 string) (*rsa.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("invalid public key PEM")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}

	return rsaPub, nil
}
