package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// hash given password
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPw), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	if tokenString == "" || tokenSecret == "" {
		return uuid.Nil, errors.New("token string or token secret is empty")
	}

	// issuer is the one issuing the JWT, usually the server.
	//
	// subject is usually who the JWT is intended for.
	//
	// audience is which resource server or API the JWT is intended for.
	// subject is within the audience.

	parser := jwt.NewParser(
		jwt.WithIssuer("chirpy"),
		jwt.WithTimeFunc(time.Now().UTC),
	)

	token, err := parser.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
			}

			return []byte(tokenSecret), nil
		})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, jwt.ErrInvalidType
	}
	if claims.Subject == "" {
		return uuid.Nil, jwt.ErrTokenInvalidSubject
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%v. invalid subject format: %w", jwt.ErrTokenInvalidSubject, err)
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	// 1. get authorization header
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	// 2. split auth header
	// auth header looks like this:
	// Authorization: <auth-scheme> <authorization-parameters>
	//
	// in this case, we are interested in the <Bearer> auth scheme
	authSplit := strings.Fields(authHeader)
	if len(authSplit) != 2 {
		return "", fmt.Errorf("invalid authorization header format: %v", authSplit)
	}
	authScheme, authToken := authSplit[0], authSplit[1]

	// 3. check auth scheme if it exists and if contains "Bearer"
	// if checks pass, continue
	if authScheme == "" || authScheme != "Bearer" {
		return "", fmt.Errorf("unexpected authorizaton scheme or is nil: %v", authScheme)
	}

	// 4. check if auth token is not nil
	if authToken == "" {
		return "", errors.New("authorization token string is nil")
	}

	return authToken, nil
}

func MakeRefreshToken() (string, error) {
	// create 32 bytes of random data and hex it
	rnd := make([]byte, 32)
	_, err := rand.Read(rnd)
	if err != nil {
		return "", fmt.Errorf("failed to read to []byte: %w", err)
	}

	hexString := hex.EncodeToString(rnd)

	return hexString, nil
}
