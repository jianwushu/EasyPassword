package auth

import (
	"context"
	"crypto/subtle"
	"easy-password-backend/config"
	"easy-password-backend/internal/apierror"
	"easy-password-backend/internal/core"
	"easy-password-backend/internal/crypto"
	"easy-password-backend/internal/email"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"
)

// AuthService 提供用户身份验证相关的服务。
type AuthService struct {
	userRepo core.UserRepository
	vcRepo   core.VerificationCodeRepository
	emailSvc email.EmailService
	cfg      *config.Config
}

// NewAuthService 创建一个新的 AuthService。
func NewAuthService(userRepo core.UserRepository, vcRepo core.VerificationCodeRepository, emailSvc email.EmailService, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		vcRepo:   vcRepo,
		emailSvc: emailSvc,
		cfg:      cfg,
	}
}

// Register 处理用户注册的业务逻辑。
func (s *AuthService) Register(ctx context.Context, username, email, masterKeyHash, masterSalt, code string) (*core.User, error) {
	slog.Info("Attempting to register new user", "username", username, "email", email)
	// 1. 验证验证码
	vc, err := s.vcRepo.Find(ctx, email)
	if err != nil {
		if err == core.ErrVerificationCodeNotFound {
			slog.Warn("Registration failed: verification code not found", "email", email)
			return nil, apierror.ErrInvalidVerificationCode
		}
		slog.Error("Registration failed: error finding verification code", "error", err)
		return nil, apierror.ErrInternalServer
	}

	if vc.Code != code {
		slog.Warn("Registration failed: invalid verification code", "email", email)
		return nil, apierror.ErrInvalidVerificationCode
	}

	if time.Now().After(vc.ExpiresAt) {
		slog.Warn("Registration failed: verification code expired", "email", email)
		return nil, apierror.ErrVerificationCodeExpired
	}

	// 2. 检查用户或邮箱是否已存在。
	_, err = s.userRepo.FindByUsername(ctx, username)
	if err == nil {
		slog.Warn("Registration failed: username already exists", "username", username)
		return nil, apierror.ErrUserOrEmailExists
	}
	if err != core.ErrUserNotFound {
		slog.Error("Registration failed: error checking username", "error", err)
		return nil, apierror.ErrInternalServer
	}

	_, err = s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		slog.Warn("Registration failed: email already exists", "email", email)
		return nil, apierror.ErrUserOrEmailExists
	}
	if err != core.ErrUserNotFound {
		slog.Error("Registration failed: error checking email", "error", err)
		return nil, apierror.ErrInternalServer
	}

	// 3. 创建一个新的用户实体。
	newUser := &core.User{
		Username:   username,
		Email:      email,
		AuthHash:   masterKeyHash, // 我们现在直接存储主密钥哈希。
		MasterSalt: []byte(masterSalt),
	}

	// 4. 将新用户保存到存储库。
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		slog.Error("Failed to create user in repository", "error", err)
		return nil, apierror.ErrInternalServer
	}

	// 5. 删除已使用的验证码
	_ = s.vcRepo.Delete(ctx, email)

	slog.Info("User registered successfully", "user_id", newUser.ID, "username", newUser.Username)
	return newUser, nil
}

// Login 处理用户登录的业务逻辑，并返回一个 JWT 和用户的主盐。
func (s *AuthService) Login(ctx context.Context, identifier, masterKeyHash string) (string, string, string, error) {
	slog.Info("Login attempt", "identifier", identifier)
	var user *core.User
	var err error
	// 1. 按标识符（用户名或邮箱）查找用户。
	if strings.Contains(identifier, "@") {
		user, err = s.userRepo.FindByEmail(ctx, identifier)
	} else {
		user, err = s.userRepo.FindByUsername(ctx, identifier)
	}

	if err != nil {
		slog.Warn("Login failed: user not found", "identifier", identifier, "error", err)
		return "", "", "", apierror.ErrInvalidCredentials
	}

	// 2. 将提供的主密钥哈希与存储的哈希进行比较。
	// 使用恒定时间比较函数来防止时序攻击。
	if len(user.AuthHash) != len(masterKeyHash) {
		slog.Warn("Login failed: invalid credentials (hash length mismatch)", "user_id", user.ID)
		return "", "", "", apierror.ErrInvalidCredentials
	}

	if subtle.ConstantTimeCompare([]byte(user.AuthHash), []byte(masterKeyHash)) == 0 {
		slog.Warn("Login failed: invalid credentials (hash mismatch)", "user_id", user.ID)
		return "", "", "", apierror.ErrInvalidCredentials
	}

	// 3. 生成 JWT。
	token, err := crypto.GenerateJWT(user.ID, s.cfg.JWTSecret, s.cfg.JWTExpiration)
	if err != nil {
		slog.Error("Failed to generate JWT for user", "user_id", user.ID, "error", err)
		return "", "", "", apierror.ErrInternalServer
	}

	slog.Info("User logged in successfully", "user_id", user.ID)
	// 4. 返回令牌和主盐、用户名。
	return token, string(user.MasterSalt), user.Username, nil
}

// GetMasterSalt 检索给定用户的主盐。
func (s *AuthService) GetMasterSalt(ctx context.Context, identifier string) (string, error) {
	var user *core.User
	var err error
	if strings.Contains(identifier, "@") {
		user, err = s.userRepo.FindByEmail(ctx, identifier)
	} else {
		user, err = s.userRepo.FindByUsername(ctx, identifier)
	}

	if err != nil {
		// 重要的是不要透露用户是否存在。
		// 返回一个通用的错误消息。
		return "", apierror.ErrInvalidCredentials
	}
	return string(user.MasterSalt), nil
}

// SendVerificationCode 生成、存储并发送一个邮件验证码。
func (s *AuthService) SendVerificationCode(ctx context.Context, emailAddr string) error {
	// 1. 检查邮箱是否已经被注册
	_, err := s.userRepo.FindByEmail(ctx, emailAddr)
	if err == nil {
		return apierror.ErrEmailExists
	}
	if err != core.ErrUserNotFound {
		return apierror.ErrInternalServer
	}

	// 2. 生成一个6位数的随机验证码
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

	// 打印验证码到控制台以供测试
	slog.Debug("Verification code generated", "email", emailAddr, "code", code)

	// 3. 创建验证码实体
	vc := &core.VerificationCode{
		Email:     emailAddr,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5分钟有效期
	}

	// 4. 存储验证码到数据库 (如果已存在则更新)
	if err := s.vcRepo.Create(ctx, vc); err != nil {
		return apierror.ErrInternalServer
	}

	// 5. 发送邮件
	// 在一个 goroutine 中发送以避免阻塞请求
	go func() {
		err := s.emailSvc.SendVerificationCodeEmail(emailAddr, code)
		if err != nil {
			slog.Error("Failed to send verification email", "recipient", emailAddr, "error", err)
		}
	}()

	return nil
}

// VerifyPasswordResetToken 验证密码重置令牌并返回用户的 MasterSalt。
// 如果令牌有效，它将返回 MasterSalt，前端可以使用它来哈希新密码。
func (s *AuthService) VerifyPasswordResetToken(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", apierror.ErrInvalidResetToken
	}
	tokenHash := crypto.HashString(token)

	user, err := s.userRepo.FindByResetPasswordToken(ctx, tokenHash)
	if err != nil {
		if err == core.ErrUserNotFound {
			return "", apierror.ErrInvalidResetToken
		}
		return "", apierror.ErrInternalServer
	}

	// 检查令牌是否已过期。
	if user.ResetPasswordTokenExpiresAt == nil || time.Now().After(*user.ResetPasswordTokenExpiresAt) {
		// 为了安全起见，清除过期的令牌。
		user.ResetPasswordToken = nil
		user.ResetPasswordTokenExpiresAt = nil
		_ = s.userRepo.Update(ctx, user)
		return "", apierror.ErrResetTokenExpired
	}

	return string(user.MasterSalt), nil
}

// RequestPasswordReset 处理密码重置请求。
func (s *AuthService) RequestPasswordReset(ctx context.Context, emailAddr string) error {
	// 1. 按邮箱查找用户。
	user, err := s.userRepo.FindByEmail(ctx, emailAddr)
	if err != nil {
		// 出于安全原因，即使找不到用户，我们也假装成功。
		// 这可以防止攻击者使用此端点来确定哪些电子邮件已注册。
		if err == core.ErrUserNotFound {
			slog.Info("Password reset requested for non-existent user", "email", emailAddr)
			return nil
		}
		slog.Error("Error finding user for password reset", "email", emailAddr, "error", err)
		return apierror.ErrInternalServer
	}

	// 2. 生成一个安全的随机令牌。
	token, err := crypto.GenerateRandomString(32)
	if err != nil {
		return apierror.ErrInternalServer
	}

	// 3. 存储令牌的哈希值和过期时间。
	// 我们存储哈希值而不是原始令牌，以增加安全性。
	tokenHash := crypto.HashString(token)
	expiresAt := time.Now().Add(30 * time.Minute) // 30分钟有效期

	user.ResetPasswordToken = &tokenHash
	user.ResetPasswordTokenExpiresAt = &expiresAt

	if err := s.userRepo.Update(ctx, user); err != nil {
		return apierror.ErrInternalServer
	}

	// 4. 发送密码重置邮件。
	// 在一个 goroutine 中发送以避免阻塞。
	go func() {
		resetLink := fmt.Sprintf("%s/reset-password/%s", s.cfg.FrontendURL, token)

		// 打印验证码到控制台以供测试
		slog.Debug("Password reset link generated", "email", emailAddr, "link", resetLink)

		err := s.emailSvc.SendPasswordResetEmail(user.Email, resetLink)
		if err != nil {
			slog.Error("Failed to send password reset email", "recipient", user.Email, "error", err)
		}
	}()

	return nil
}

// ResetPassword 使用有效的重置令牌重置用户的密码。
func (s *AuthService) ResetPassword(ctx context.Context, token, newMasterKeyHash, newMasterSalt string) error {
	slog.Info("Attempting to reset password")
	// 1. 验证令牌。
	if token == "" {
		slog.Warn("Password reset failed: no token provided")
		return apierror.ErrInvalidResetToken
	}
	tokenHash := crypto.HashString(token)

	user, err := s.userRepo.FindByResetPasswordToken(ctx, tokenHash)
	if err != nil {
		if err == core.ErrUserNotFound {
			slog.Warn("Password reset failed: invalid token")
			return apierror.ErrInvalidResetToken
		}
		slog.Error("Error finding user by reset token", "error", err)
		return apierror.ErrInternalServer
	}

	// 2. 检查令牌是否已过期。
	if user.ResetPasswordTokenExpiresAt == nil || time.Now().After(*user.ResetPasswordTokenExpiresAt) {
		slog.Warn("Password reset failed: token expired", "user_id", user.ID)
		// 为了安全起见，清除过期的令牌。
		user.ResetPasswordToken = nil
		user.ResetPasswordTokenExpiresAt = nil
		_ = s.userRepo.Update(ctx, user)
		return apierror.ErrResetTokenExpired
	}

	// 3. 更新用户的 AuthHash 和 MasterSalt。
	user.AuthHash = newMasterKeyHash
	user.MasterSalt = []byte(newMasterSalt)

	// 4. 使重置令牌失效。
	user.ResetPasswordToken = nil
	user.ResetPasswordTokenExpiresAt = nil

	if err := s.userRepo.Update(ctx, user); err != nil {
		slog.Error("Failed to update user password", "user_id", user.ID, "error", err)
		return apierror.ErrInternalServer
	}

	slog.Info("Password reset successfully", "user_id", user.ID)
	return nil
}
