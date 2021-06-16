package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fr13n8/todo-app"
	"github.com/fr13n8/todo-app/pkg/repository"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

var (
	cost       = viper.GetInt("password.cost")
	signingKey = viper.GetString("jwt.signingKey")
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.SignUpInput) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) CreateSession(input todo.Session) error {
	return s.repo.CreateSession(input)
}

func (s *AuthService) SignInUser(username, password, userAgent string) ([]string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	tokens, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	session := todo.Session{
		UserId:       user.Id,
		RefreshToken: tokens[1],
		UUID:         uuid,
		UserAgent:    userAgent,
	}
	if err := s.repo.CreateSession(session); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthService) GenerateToken(user todo.User) ([]string, error) {

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.UserName,
		ExpiresAt: now.Add(12 * time.Hour).Unix(),
		IssuedAt:  now.Unix(),
		Id:        strconv.Itoa(user.Id),
	})
	accessToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return nil, err
	}

	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Issuer:    user.UserName,
		ExpiresAt: now.Add(1 * time.Minute).Unix(),
		IssuedAt:  now.Unix(),
	})
	refreshToken, err := rToken.SignedString([]byte(signingKey))
	if err != nil {
		return nil, err
	}

	return []string{accessToken, refreshToken}, nil
}

func (s *AuthService) ParseToken(accessToken string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.New("token claims are not found")
	}

	return claims, nil
}

func (s *AuthService) RefreshToken(token string) ([]string, error) {
	claims, err := s.ParseToken(token)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUser(claims.Issuer)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("refresh token expired")
	}

	tokens, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func generatePasswordHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hash)
}
