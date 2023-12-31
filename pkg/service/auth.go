package service

import (
	"copySys/models"
	"copySys/pkg/logger"
	"copySys/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	salt = "ajsdijaskdasl122312klsdjka"
	// будем использовать для рассшифровки токена
	signingKey = "kajsdljaskdja332$#"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId   int    `json:"user_id"`
	Role     string `json:"role"`
	UserName string `json:"user_name"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)

	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) GenerateToken(username, password, role string) (string, error) {

	password = generatePasswordHash(password)
	user, err := s.repo.GetUser(username, password, role)

	if err != nil {
		logger.Error.Println("GetUser func: ", err.Error())
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(120 * time.Hour).Unix(), // время жизни токена
			IssuedAt:  time.Now().Unix(),                      // когда был создан токен
		},
		user.ID,
		user.Role,
		user.UserName,
	})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, string, string, error) {
	// ключ подпись или ошибку

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// метод проверки токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Error.Println("invalid signing method")
			return nil, errors.New("invalid signing method")

		}

		return []byte(signingKey), nil
	})

	if err != nil {
		logger.Error.Println(err)
		return 0, "", "", err
	}

	claims, ok := token.Claims.(*tokenClaims)

	if !ok {
		logger.Error.Println("token claims are not of type *tokenClaims ")
		return 0, "", "", errors.New("token claims are not of type *tokenClaims ")
	}

	if time.Now().Unix() > claims.ExpiresAt {
		logger.Error.Println("token expired")
		return 0, "", "", errors.New("token expired")
	}

	return claims.UserId, claims.Role, claims.UserName, nil
}
