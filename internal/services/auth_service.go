package services

import (
	"context"
	"errors"
	"os"
	"time"

	"altoai_mvp/internal/models"
	"altoai_mvp/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, dto models.LoginDTO) (string, *models.User, error)
	Register(ctx context.Context, dto models.CreateUserDTO) error
	VerifyEmail(ctx context.Context, dto models.VerifyEmailDTO) (string, *models.User, error)
	ResendVerificationCode(ctx context.Context, dto models.ResendVerificationDTO) error
	ForgotPassword(ctx context.Context, dto models.ForgotPasswordDTO) error
	ResetPassword(ctx context.Context, dto models.ResetPasswordDTO) error
	RefreshToken(ctx context.Context, tokenString string) (string, error)
}

type authService struct {
	userRepo repository.UserRepo
	emailSvc EmailService
}

func NewAuthService(userRepo repository.UserRepo) AuthService {
	return &authService{
		userRepo: userRepo,
		emailSvc: NewEmailService(),
	}
}

func (s *authService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *authService) comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *authService) generateToken(user models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	claims := jwt.MapClaims{
		"email":   user.Email,
		"name":    user.Name,
		"picture": "",
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     "altoai_mvp",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *authService) Login(ctx context.Context, dto models.LoginDTO) (string, *models.User, error) {
	user, err := s.userRepo.GetByEmail(dto.Email)
	if err != nil {
		if err == repository.ErrNotFound {
			return "", nil, errors.New("invalid email or password")
		}
		return "", nil, err
	}

	// Check if email is verified
	if !user.EmailVerified {
		return "", nil, errors.New("email not verified. Please check your email for verification code")
	}

	// Check if user has a password (OAuth users might not have one)
	if user.Password == "" {
		return "", nil, errors.New("invalid email or password")
	}

	// Verify password
	if err := s.comparePassword(user.Password, dto.Password); err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

func (s *authService) Register(ctx context.Context, dto models.CreateUserDTO) error {
	// Password is required for registration
	if dto.Password == "" {
		return errors.New("password is required")
	}

	// Check if user already exists
	_, err := s.userRepo.GetByEmail(dto.Email)
	if err == nil {
		return errors.New("user with this email already exists")
	}
	if err != repository.ErrNotFound {
		return err
	}

	// Hash password
	passwordHash, err := s.hashPassword(dto.Password)
	if err != nil {
		return err
	}

	// Create user
	user, err := s.userRepo.Create(dto.Email, dto.Name, passwordHash)
	if err != nil {
		return err
	}

	// Generate verification code
	code, err := s.emailSvc.GenerateCode()
	if err != nil {
		return err
	}

	// Set verification code (expires in 15 minutes) - use UTC
	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	if err := s.userRepo.SetVerificationCode(user.Email, code, expiresAt); err != nil {
		return err
	}

	// Send verification email
	if err := s.emailSvc.SendVerificationCode(user.Email, user.Name, code); err != nil {
		return errors.New("failed to send verification email")
	}

	return nil
}

func (s *authService) VerifyEmail(ctx context.Context, dto models.VerifyEmailDTO) (string, *models.User, error) {
	// Verify the code
	if err := s.userRepo.VerifyEmail(dto.Email, dto.Code); err != nil {
		return "", nil, err
	}

	// Get the verified user
	user, err := s.userRepo.GetByEmail(dto.Email)
	if err != nil {
		return "", nil, err
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

func (s *authService) ResendVerificationCode(ctx context.Context, dto models.ResendVerificationDTO) error {
	// Check if user exists
	user, err := s.userRepo.GetByEmail(dto.Email)
	if err != nil {
		if err == repository.ErrNotFound {
			// Don't reveal if user exists or not for security
			return nil
		}
		return err
	}

	// Check if email is already verified
	if user.EmailVerified {
		return errors.New("email is already verified")
	}

	// Generate new verification code
	code, err := s.emailSvc.GenerateCode()
	if err != nil {
		return err
	}

	// Set verification code (expires in 15 minutes) - use UTC
	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	if err := s.userRepo.SetVerificationCode(user.Email, code, expiresAt); err != nil {
		return err
	}

	// Send verification email
	if err := s.emailSvc.SendVerificationCode(user.Email, user.Name, code); err != nil {
		return errors.New("failed to send verification email")
	}

	return nil
}

func (s *authService) ForgotPassword(ctx context.Context, dto models.ForgotPasswordDTO) error {
	// Check if user exists
	user, err := s.userRepo.GetByEmail(dto.Email)
	if err != nil {
		if err == repository.ErrNotFound {
			// Don't reveal if user exists or not for security
			return nil
		}
		return err
	}

	// Generate reset code
	code, err := s.emailSvc.GenerateCode()
	if err != nil {
		return err
	}

	// Set reset code (expires in 15 minutes) - use UTC
	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	if err := s.userRepo.SetResetCode(user.Email, code, expiresAt); err != nil {
		return err
	}

	// Send reset email
	if err := s.emailSvc.SendPasswordResetCode(user.Email, user.Name, code); err != nil {
		return errors.New("failed to send reset email")
	}

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, dto models.ResetPasswordDTO) error {
	// Verify reset code and update password
	passwordHash, err := s.hashPassword(dto.Password)
	if err != nil {
		return err
	}

	if err := s.userRepo.ResetPassword(dto.Email, dto.Code, passwordHash); err != nil {
		return err
	}

	return nil
}

func (s *authService) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// Get user email from claims
	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// Get user from database
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Generate new token
	return s.generateToken(user)
}
