package authentication

import (
	"errors"
	"time"

	"github.com/FranzSinaga/blogcms/internal/domain"
	"github.com/FranzSinaga/blogcms/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo  RepositoryInterface
	jwtConfig config.JWTConfig
}

func NewAuthService(userRepo RepositoryInterface, jwtConfig config.JWTConfig) *Service {
	return &Service{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
	}
}

func (s *Service) Register(req *domain.RegisterRequest) (string, *domain.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	user := &domain.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     "user",
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return "", nil, err
	}

	tokenString, err := s.generateToken(user)
	return tokenString, &domain.UserResponse{
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *Service) Login(req *domain.LoginRequest) (string, *domain.UserResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	tokenString, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return tokenString, &domain.UserResponse{
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *Service) generateToken(user *domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"name":    user.Name,
		"exp":     time.Now().Add(s.jwtConfig.ExpiresIn).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     "cms-api",
	})
	return token.SignedString([]byte(s.jwtConfig.Secret))
}
