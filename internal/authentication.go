package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Tokenize(id int, secret_key string, expiresin time.Duration) (string, error) {
	secret_key_byte := []byte(secret_key)
	fmt.Println(expiresin)
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresin)),
		Subject:   strconv.Itoa(id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret_key_byte)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func BearerHeader(headers http.Header) (string, error) {

	token := headers.Get("Authorization")

	if token == "" {
		return "", errors.New("Auth header not found")
	}

	splitToken := strings.Split(token, " ")

	if len(splitToken) < 2 || splitToken[0] != "Bearer" {
		return "", errors.New("Auth Header not what expected")
	}

	return splitToken[1], nil
}

func ValidateToken(tokenstring, tokenSecret string) (string, error) {
	type customClaims struct {
		jwt.RegisteredClaims
	}
	token, err := jwt.ParseWithClaims(tokenstring, &customClaims{}, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })

	if err != nil {
		return "", errors.New(err.Error()) //"jwt couldn't be parsed"
	}

	userId, err := token.Claims.GetSubject()

	if err != nil {
		return "", errors.New("User id couldn't be extracted")
	}

	return userId, nil
}
