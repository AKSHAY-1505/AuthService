package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"log"
	"strings"

	"github.com/AKSHAY-1505/auth-service/initializers"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

var (
	privateKey jwk.Key
	publicKey  jwk.Key
	jwks       jwk.Set
)

func InitJWTSecretKeys() {
	var err error

	privateKeyBase64 := initializers.AppConfig.JWTPrivateKey
	if privateKeyBase64 == "" {
		panic("private key not configured in environment variables")
	}

	publicKeyBase64 := initializers.AppConfig.JWTPublicKey
	if publicKeyBase64 == "" {
		panic("public key not configured in environment variables")
	}

	privateKey, err = loadJWKFromBase64(privateKeyBase64)
	if err != nil {
		log.Fatalf("failed to parse JWT private key: %v", err)
	}

	publicKey, err = loadJWKFromBase64(publicKeyBase64)
	if err != nil {
		log.Fatalf("failed to parse JWT public key: %v", err)
	}

	jwks, err = buildJWKS(publicKey)
	if err != nil {
		log.Fatalf("failed to create jwks: %v", err)
	}
}

func GetPrivateKey() *rsa.PrivateKey {
	var key rsa.PrivateKey
	_ = privateKey.Raw(&key)
	return &key
}

func GetPublicKey() *rsa.PublicKey {
	var key rsa.PublicKey
	_ = privateKey.Raw(&key)
	return &key
}

func GetJWKS() jwk.Set {
	return jwks
}

func loadJWKFromBase64(b64 string) (jwk.Key, error) {
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}

	key, err := jwk.ParseKey(decoded, jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}

	return key, nil
}

// BuildJWKS returns a jwk.Set containing the public key
func buildJWKS(publicKey jwk.Key) (jwk.Set, error) {
	// Ensure required fields are set
	// You SHOULD already be setting kid during key generation
	if _, ok := publicKey.Get(jwk.KeyIDKey); !ok {
		_ = publicKey.Set(jwk.KeyIDKey, "default-kid")
	}

	// Set recommended metadata
	_ = publicKey.Set(jwk.AlgorithmKey, "RS256")
	_ = publicKey.Set(jwk.KeyUsageKey, "sig")

	set := jwk.NewSet()
	if err := set.AddKey(publicKey); err != nil {
		return nil, err
	}

	return set, nil
}

func ExtractAuthHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("authorization header is missing")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	return tokenString, nil
}
