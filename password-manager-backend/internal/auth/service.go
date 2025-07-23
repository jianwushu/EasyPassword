package auth

import (
	"context"
	"easy-password-backend/config"
	"easy-password-backend/internal/apierror"
	"easy-password-backend/internal/core"
	"easy-password-backend/internal/crypto"
)

// AuthService 提供用户身份验证相关的服务。
type AuthService struct {
	userRepo core.UserRepository
	cfg      *config.Config
}

// NewAuthService 创建一个新的 AuthService。
func NewAuthService(userRepo core.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

// Register 处理用户注册的业务逻辑。
func (s *AuthService) Register(ctx context.Context, username, masterKeyHash, masterSalt string) (*core.User, error) {
	// 1. 检查用户是否已存在。
	_, err := s.userRepo.FindByUsername(ctx, username)
	if err == nil {
		return nil, apierror.ErrUsernameExists
	}

	// 2. 创建一个新的用户实体。
	newUser := &core.User{
		Username:   username,
		AuthHash:   masterKeyHash, // 我们现在直接存储主密钥哈希。
		MasterSalt: []byte(masterSalt),
	}

	// 3. 将新用户保存到存储库。
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, apierror.ErrInternalServer
	}

	return newUser, nil
}

// Login 处理用户登录的业务逻辑，并返回一个 JWT 和用户的主盐。
func (s *AuthService) Login(ctx context.Context, username, masterKeyHash string) (string, string, error) {
	// 1. 按用户名查找用户。
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", "", apierror.ErrInvalidCredentials
	}

	// 2. 将提供的主密钥哈希与存储的哈希进行比较。
	// 注意：这是一个直接的字符串比较。在实际场景中，
	// 您应该使用恒定时间比较函数来防止时序攻击。
	if user.AuthHash != masterKeyHash {
		return "", "", apierror.ErrInvalidCredentials
	}

	// 3. 生成 JWT。
	token, err := crypto.GenerateJWT(user.ID, s.cfg.JWTSecret, s.cfg.JWTExpiration)
	if err != nil {
		return "", "", apierror.ErrInternalServer
	}

	// 4. 返回令牌和主盐。
	return token, string(user.MasterSalt), nil
}

// GetMasterSalt 检索给定用户的主盐。
func (s *AuthService) GetMasterSalt(ctx context.Context, username string) (string, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		// 重要的是不要透露用户是否存在。
		// 返回一个通用的错误消息。
		return "", apierror.ErrInvalidCredentials
	}
	return string(user.MasterSalt), nil
}