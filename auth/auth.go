package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/AKSHAY-1505/auth-service/initializers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	_ = publicKey.Raw(&key)
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

func CreateJWT(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(GetPrivateKey())
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseClaims(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Enforce expected signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return GetPublicKey(), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Optional: manually verify exp if you want stricter control
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token expired")
			}
		}

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
