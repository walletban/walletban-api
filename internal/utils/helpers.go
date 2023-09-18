package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"net/http"
	"time"
)

func GenerateOauthHash(url string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	claims["url"] = url
	tokenString, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func ValidateToken(encodedToken string) (string, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{},
		error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token %s", token.Header["alg"])
		}
		return []byte(JwtSecret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["url"].(string), nil
	} else {
		return "", errors.New("invalid claims")
	}
}

func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Error("Error generating random string: %v\n", err)
		return "", err
	}
	for i := 0; i < length; i++ {
		randomBytes[i] = charset[int(randomBytes[i])%len(charset)]
	}
	return string(randomBytes), nil
}

func MakePOSTRequest(url string, data []byte) (map[string]interface{}, error) {

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new POST request with the JSON payload
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response body into a map
	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return nil, err
	}

	return responseMap, nil
}
