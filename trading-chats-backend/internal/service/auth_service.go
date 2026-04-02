package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"trading-chats-backend/internal/config"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/repository"
	"trading-chats-backend/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       *repository.AuthRepository
	jwtConfig  config.JWTConfig
	refreshTTL time.Duration
}

func NewAuthService(repo *repository.AuthRepository, jwtConfig config.JWTConfig) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtConfig:  jwtConfig,
		refreshTTL: 7 * 24 * time.Hour,
	}
}

func (s *AuthService) EnsureBootstrapData(ctx context.Context) error {
	if err := s.repo.EnsureIndexes(ctx); err != nil {
		return err
	}

	seedTenants := []models.Tenant{
		{ID: "tenant-system", Name: "System Tenant", Code: "system", Type: models.TenantTypeSystem, Status: "active"},
		{ID: "tenant-alpha", Name: "Tenant Alpha", Code: "tenant_alpha", Type: models.TenantTypeBiz, Status: "active"},
		{ID: "tenant-beta", Name: "Tenant Beta", Code: "tenant_beta", Type: models.TenantTypeBiz, Status: "active"},
	}
	for _, tenant := range seedTenants {
		t := tenant
		if err := s.repo.UpsertTenant(ctx, &t); err != nil {
			return err
		}
	}

	adminPassword := getSeedPassword("SEED_ADMIN_PASSWORD", "Admin@123456")
	alphaPassword := getSeedPassword("SEED_TENANT_ALPHA_PASSWORD", "TenantAlpha@123456")
	betaPassword := getSeedPassword("SEED_TENANT_BETA_PASSWORD", "TenantBeta@123456")

	adminHash, err := HashPassword(adminPassword)
	if err != nil {
		return err
	}
	alphaHash, err := HashPassword(alphaPassword)
	if err != nil {
		return err
	}
	betaHash, err := HashPassword(betaPassword)
	if err != nil {
		return err
	}

	seedUsers := []models.User{
		{ID: "user-admin", TenantID: "tenant-system", Username: "admin", PasswordHash: adminHash, DisplayName: "System Admin", Role: models.RoleAdmin, Status: "active"},
		{ID: "user-tenant-alpha", TenantID: "tenant-alpha", Username: "tenant_alpha", PasswordHash: alphaHash, DisplayName: "Tenant Alpha Admin", Role: models.RoleTenant, Status: "active"},
		{ID: "user-tenant-beta", TenantID: "tenant-beta", Username: "tenant_beta", PasswordHash: betaHash, DisplayName: "Tenant Beta Admin", Role: models.RoleTenant, Status: "active"},
	}
	for _, user := range seedUsers {
		u := user
		if err := s.repo.UpsertUser(ctx, &u); err != nil {
			return err
		}
	}

	return nil
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.LoginResponse, error) {
	if req == nil {
		return nil, errors.New("invalid register request")
	}

	tenantCode := strings.TrimSpace(req.TenantCode)
	tenantName := strings.TrimSpace(req.TenantName)
	username := strings.TrimSpace(req.Username)
	displayName := strings.TrimSpace(req.DisplayName)
	password := strings.TrimSpace(req.Password)

	if tenantCode == "" || tenantName == "" || username == "" || displayName == "" || password == "" {
		return nil, errors.New("tenant_name, tenant_code, username, display_name and password are required")
	}

	if _, err := s.repo.GetTenantByCode(ctx, tenantCode); err == nil {
		return nil, errors.New("tenant code already exists")
	}

	if _, err := s.repo.GetUserByUsername(ctx, username); err == nil {
		return nil, errors.New("username already exists")
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	tenant := &models.Tenant{
		ID:        tenantCode,
		Name:      tenantName,
		Code:      tenantCode,
		Type:      models.TenantTypeBiz,
		Status:    "active",
		CreatedAt: utils.NowString(),
		UpdatedAt: utils.NowString(),
	}
	if err := s.repo.CreateTenant(ctx, tenant); err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           "user-" + tenantCode + "-" + username,
		TenantID:     tenant.ID,
		Username:     username,
		PasswordHash: passwordHash,
		DisplayName:  displayName,
		Role:         models.RoleTenant,
		Status:       "active",
		CreatedAt:    utils.NowString(),
		UpdatedAt:    utils.NowString(),
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return s.issueSession(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	if req == nil {
		return nil, errors.New("invalid login request")
	}

	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	if user.Status != "active" {
		return nil, errors.New("user is disabled")
	}
	if err := ComparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return s.issueSession(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.LoginResponse, error) {
	if req == nil || req.RefreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	session, err := s.repo.GetSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	if session.ExpiresAt.Before(time.Now()) {
		_ = s.repo.DeleteSessionByRefreshToken(ctx, req.RefreshToken)
		return nil, errors.New("refresh token expired")
	}

	user, err := s.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.DeleteSessionByRefreshToken(ctx, req.RefreshToken); err != nil {
		return nil, err
	}

	return s.issueSession(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, req *models.LogoutRequest) error {
	if req == nil || req.RefreshToken == "" {
		return errors.New("refresh token is required")
	}
	return s.repo.DeleteSessionByRefreshToken(ctx, req.RefreshToken)
}

func (s *AuthService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	if req == nil {
		return errors.New("invalid reset password request")
	}

	username := strings.TrimSpace(req.Username)
	newPassword := strings.TrimSpace(req.NewPassword)
	if username == "" || newPassword == "" {
		return errors.New("username and new_password are required")
	}

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return errors.New("user not found")
	}

	passwordHash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateUserPassword(ctx, user.ID, passwordHash); err != nil {
		return err
	}
	if err := s.repo.DeleteSessionsByUserID(ctx, user.ID); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ValidateAccessToken(ctx context.Context, accessToken string) (*models.AuthContext, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	session, err := s.repo.GetSessionByAccessToken(ctx, accessToken)
	if err != nil {
		return nil, errors.New("session not found")
	}
	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session expired")
	}

	userID, _ := claims["user_id"].(string)
	tenantID, _ := claims["tenant_id"].(string)
	roleText, _ := claims["role"].(string)
	username, _ := claims["username"].(string)
	sessionID, _ := claims["session_id"].(string)

	return &models.AuthContext{UserID: userID, TenantID: tenantID, Role: models.Role(roleText), Username: username, SessionID: sessionID}, nil
}

func (s *AuthService) issueSession(ctx context.Context, user *models.User) (*models.LoginResponse, error) {
	if err := s.repo.DeleteExpiredSessions(ctx, time.Now()); err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(s.jwtConfig.Expiration)
	sessionID, err := generateRandomToken(16)
	if err != nil {
		return nil, err
	}
	accessToken, err := s.generateAccessToken(user, sessionID, expiresAt)
	if err != nil {
		return nil, err
	}
	refreshToken, err := generateRandomToken(32)
	if err != nil {
		return nil, err
	}

	session := &models.UserSession{ID: sessionID, UserID: user.ID, TenantID: user.TenantID, Role: user.Role, AccessToken: accessToken, RefreshToken: refreshToken, ExpiresAt: time.Now().Add(s.refreshTTL), CreatedAt: utils.NowString()}
	if err := s.repo.SaveSession(ctx, session); err != nil {
		return nil, err
	}
	_ = s.repo.UpdateUserLastLogin(ctx, user.ID)

	cleanUser := *user
	cleanUser.PasswordHash = ""

	return &models.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken, TokenType: "Bearer", ExpiresAt: expiresAt, User: cleanUser}, nil
}

func (s *AuthService) generateAccessToken(user *models.User, sessionID string, expiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{"user_id": user.ID, "tenant_id": user.TenantID, "role": string(user.Role), "username": user.Username, "session_id": sessionID, "exp": expiresAt.Unix(), "iat": time.Now().Unix()}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtConfig.Secret))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func generateRandomToken(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

func getSeedPassword(envKey, fallback string) string {
	if value := os.Getenv(envKey); value != "" {
		return value
	}
	return fallback
}
