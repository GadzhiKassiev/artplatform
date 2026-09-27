package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"artplatform/backend/internal/pkg/jwt"
	"artplatform/backend/internal/service/auth/model"
	"artplatform/backend/internal/service/auth/storage"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidRole        = errors.New("invalid role")
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
	if input.Role != string(model.RoleStudent) && input.Role != string(model.RoleTeacher) {
		return AuthResult{}, ErrInvalidRole
	}

	_, err := s.storage.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return AuthResult{}, ErrEmailTaken
	}
	if !errors.Is(err, storage.ErrUserNotFound) {
		return AuthResult{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
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
		return AuthResult{}, err
	}

	token, err := s.tokenizer.Generate(created.ID, string(created.Role))
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{UserID: created.ID, Token: token}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (AuthResult, error) {
	user, err := s.storage.GetUserByEmail(ctx, email)
	if errors.Is(err, storage.ErrUserNotFound) {
		return AuthResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return AuthResult{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	token, err := s.tokenizer.Generate(user.ID, string(user.Role))
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{UserID: user.ID, Token: token}, nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (userID, role string, err error) {
	claims, err := s.tokenizer.Parse(token)
	if err != nil {
		return "", "", err
	}
	return claims.UserID, claims.Role, nil
}

func (s *Service) GetUserByID(ctx context.Context, id model.UserID) (model.User, error) {
	return s.storage.GetUserByID(ctx, id)
}
