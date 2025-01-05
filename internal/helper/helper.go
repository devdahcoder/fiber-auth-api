package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secrete-key")

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func VerifyPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func CreateToken(userEmail string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userEmail": userEmail,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtKey)
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func ExtractToken(tokenString string) (string, error) {
	claims, err := VerifyToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims["userId"].(string), nil
}

func ExtractTokenFromHeader(header string) (string, error) {
	tokenString := header[7:]
	return ExtractToken(tokenString)
}

func ExtractTokenFromCookie(cookie string) (string, error) {
	tokenString := cookie[7:]
	return ExtractToken(tokenString)
}