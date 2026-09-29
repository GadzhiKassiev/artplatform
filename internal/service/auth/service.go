package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"artplatform/backend/internal/jwt"
	"artplatform/backend/internal/logging"
	autherr "artplatform/backend/internal/service/auth/err"
	"artplatform/backend/internal/service/auth/model"
	"artplatform/backend/internal/service/auth/storage"
)

type Service struct {
	storage   storage.Storage
	tokenizer *jwt.Tokenizer
}

func New(storage storage.Storage, tokenizer *jwt.Tokenizer) *Service {
	return &Service{
		storage:   storage,
		tokenizer: tokenizer,
	}
}

type RegisterInput struct {
	Email    string
	Password string
	Role     string
}

type AuthResult struct {
	UserID string
	Token  string
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (AuthResult, error) {
	logging.ContextInfo(ctx, "registering user",
		logging.NewKV("email", input.Email),
		logging.NewKV("role", input.Role),
	)

	if input.Role != string(model.RoleStudent) && input.Role != string(model.RoleTeacher) {
		logging.ContextWarn(ctx, "invalid role", logging.NewKV("role", input.Role))
		return AuthResult{}, autherr.ErrInvalidRole
	}

	_, err := s.storage.GetUserByEmail(ctx, input.Email)
	if err == nil {
		logging.ContextWarn(ctx, "email already taken", logging.NewKV("email", input.Email))
		return AuthResult{}, autherr.ErrEmailTaken
	}
	if !errors.Is(err, storage.ErrUserNotFound) {
		logging.ContextErrorE(ctx, "failed to check email", err)
		return AuthResult{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to hash password", err)
		return AuthResult{}, err
	}

	user := model.User{
		ID:           uuid.NewString(),
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         model.Role(input.Role),
	}

	created, err := s.storage.CreateUser(ctx, user)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to create user", err)
		return AuthResult{}, err
	}

	token, err := s.tokenizer.Generate(created.ID, string(created.Role))
	if err != nil {
		logging.ContextErrorE(ctx, "failed to generate token", err)
		return AuthResult{}, err
	}

	logging.ContextInfo(ctx, "user registered", logging.NewKV("userID", created.ID))
	return AuthResult{UserID: created.ID, Token: token}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (AuthResult, error) {
	logging.ContextInfo(ctx, "login attempt", logging.NewKV("email", email))

	user, err := s.storage.GetUserByEmail(ctx, email)
	if errors.Is(err, storage.ErrUserNotFound) {
		logging.ContextWarn(ctx, "user not found", logging.NewKV("email", email))
		return AuthResult{}, autherr.ErrInvalidCredentials
	}
	if err != nil {
		logging.ContextErrorE(ctx, "failed to get user by email", err)
		return AuthResult{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		logging.ContextWarn(ctx, "invalid password", logging.NewKV("email", email))
		return AuthResult{}, autherr.ErrInvalidCredentials
	}

	token, err := s.tokenizer.Generate(user.ID, string(user.Role))
	if err != nil {
		logging.ContextErrorE(ctx, "failed to generate token", err)
		return AuthResult{}, err
	}

	logging.ContextInfo(ctx, "user logged in", logging.NewKV("userID", user.ID))
	return AuthResult{UserID: user.ID, Token: token}, nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (userID, role string, err error) {
	claims, err := s.tokenizer.Parse(token)
	if err != nil {
		logging.ContextWarn(ctx, "invalid token")
		return "", "", err
	}
	return claims.UserID, claims.Role, nil
}

func (s *Service) GetUserByID(ctx context.Context, id model.UserID) (model.User, error) {
	return s.storage.GetUserByID(ctx, id)
}
