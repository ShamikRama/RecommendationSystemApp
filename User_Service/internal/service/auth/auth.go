package auth

import (
	"User_Service/internal/converter"
	"User_Service/internal/erorrs"
	"User_Service/internal/logger"
	"User_Service/internal/repository"
	"User_Service/internal/service/model"
	"User_Service/internal/utils"
	"User_Service/pkg/kafka"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type ServiceAuth struct {
	repo   repository.Auth
	logger *logger.Logger
	kafka  kafka.Kafka
}

func NewAuthService(repo repository.Auth, logger *logger.Logger, kafka kafka.Kafka) *ServiceAuth {
	return &ServiceAuth{
		repo:   repo,
		logger: logger,
		kafka:  kafka,
	}
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

// надо сделать не здесь, а в env, перенесу если буду успевать
const (
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 7 * time.Hour
)

func (r *ServiceAuth) CreateUser(ctx context.Context, login model.Login) (int, error) {
	if len(login.Email) <= 0 || len(login.Password) <= 0 {
		r.logger.Error("Invalid input")
		zap.String("operations", "service.Auth.CreateUser")
		return 0, erorrs.ErrEmptyInput
	}

	user := converter.FromLoginToUser(login)
	user.Password = utils.GeneratePasswordHash(user.Password)

	userId, err := r.repo.Create(ctx, user)
	if err != nil {
		r.logger.Error("Failed to create user",
			zap.String("operation", "service.Auth.CreateUser"),
			zap.Error(err))

		if errors.Is(err, erorrs.ErrUniqueConstraint) {
			return 0, erorrs.ErrEmailAlreadyExists
		}

		return 0, err
	}

	msg := &kafka.KafkaMessage{
		Id:    userId,
		Event: "user_create",
		Data: map[string]interface{}{
			"email": login.Email,
		},
	}

	key := strconv.Itoa(userId)

	err = r.kafka.SendMessage(key, msg)
	if err != nil {
		r.logger.Error("Error sending msg to kafka",
			zap.String("operation", "sql.ClientRepository.GetUserEvents"),
			zap.Error(err))
		return 0, err
	}

	r.logger.Info("Success sending msg to kafka")

	return userId, nil
}

func (r *ServiceAuth) GenerateJwtToken(ctx context.Context, email string, password string) (string, error) {
	user, err := r.repo.GetUser(ctx, email, utils.GeneratePasswordHash(password))
	if err != nil {
		r.logger.Error("Error getting user",
			zap.String("operation", "sql.ClientRepository.GenerateJwtToken"),
			zap.Error(err))

		if errors.Is(err, erorrs.ErrNoRows) {
			r.logger.Error("User not found",
				zap.String("operation", "sql.ClientRepository.GenerateJwtToken"),
				zap.Error(err))
			return "", erorrs.ErrUserNotFound
		}
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		}, user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (r *ServiceAuth) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				r.logger.Error("Invalid signing method",
					zap.String("operation", "sql.ClientRepository.ParseToken"),
					zap.Bool("res", ok))

				return nil, errors.New("invalid signing method")
			}

			return []byte(signingKey), nil
		})
	if err != nil {
		r.logger.Error("Can not parse token",
			zap.String("operation", "sql.ClientRepository.ParseToken"),
			zap.Error(err))

		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		r.logger.Error("Wrong types of claims",
			zap.String("operation", "sql.ClientRepository.ParseToken"),
			zap.Bool("res", ok))

		return 0, err
	}

	return claims.UserId, nil
}
