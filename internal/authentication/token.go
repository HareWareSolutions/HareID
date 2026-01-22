package authentication

import (
	"HareID/config"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Cria um token
func CreateToken(Google_Subscription string, userID uint64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissions["User_ID"] = userID
	permissions["Google_Subscription"] = Google_Subscription

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)

	return token.SignedString([]byte(config.SecretKey))
}

// Captura o token
func GetToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

// Verifica se o token da requisição bate com chave de validação de segurança do API
func ValidateToken(r *http.Request) error {
	tokenString := GetToken(r)
	token, err := jwt.Parse(tokenString, validationKey)

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("Invalid token!")
}

// Cria a chave de validação com base na SecretKey da API
func validationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Wrong signing method! %v", token.Header["alg"])
	}

	return []byte(config.SecretKey), nil
}

// Captura o ID do usuário inserido no Token da requisição
func GetTokenUserID(r *http.Request) (string, error) {
	tokenString := GetToken(r)
	token, err := jwt.Parse(tokenString, validationKey)
	if err != nil {
		return "", err
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		val, ok := permissions["User_ID"]

		if !ok {
			return "", errors.New("chave 'User_ID' não encontrada no token")
		}

		return fmt.Sprint(val), nil
	}

	return "", errors.New("invalid token")
}

func GetTokenGoogle_Subscription(r *http.Request) (string, error) {
	tokenString := GetToken(r)
	token, err := jwt.Parse(tokenString, validationKey)
	if err != nil {
		return "", err
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		val, ok := permissions["Google_Subscription"]
		if !ok {
			return "", errors.New("chave 'Google_Subscription' não encontrada no token")
		}

		return fmt.Sprint(val), nil
	}

	return "", errors.New("invalid token")
}
