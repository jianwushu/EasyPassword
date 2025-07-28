package auth

import (
	"context"
	"easy-password-backend/config"
	"easy-password-backend/internal/apierror"
	"easy-password-backend/internal/core"
	"easy-password-backend/internal/crypto"
	"easy-password-backend/internal/email"
	"fmt"
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
	// 1. 验证验证码
	vc, err := s.vcRepo.Find(ctx, email)
	if err != nil {
		if err == core.ErrVerificationCodeNotFound {
			return nil, apierror.ErrInvalidVerificationCode
		}
		return nil, apierror.ErrInternalServer
	}

	if vc.Code != code {
		return nil, apierror.ErrInvalidVerificationCode
	}

	if time.Now().After(vc.ExpiresAt) {
		return nil, apierror.ErrVerificationCodeExpired
	}

	// 2. 检查用户或邮箱是否已存在。
	_, err = s.userRepo.FindByUsername(ctx, username)
	if err == nil {
		return nil, apierror.ErrUserOrEmailExists
	}
	if err != core.ErrUserNotFound {
		return nil, apierror.ErrInternalServer
	}

	_, err = s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, apierror.ErrUserOrEmailExists
	}
	if err != core.ErrUserNotFound {
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
		return nil, apierror.ErrInternalServer
	}

	// 5. 删除已使用的验证码
	_ = s.vcRepo.Delete(ctx, email)

	return newUser, nil
}

// Login 处理用户登录的业务逻辑，并返回一个 JWT 和用户的主盐。
func (s *AuthService) Login(ctx context.Context, identifier, masterKeyHash string) (string, string, string, error) {
	var user *core.User
	var err error
	// 1. 按标识符（用户名或邮箱）查找用户。
	if strings.Contains(identifier, "@") {
		user, err = s.userRepo.FindByEmail(ctx, identifier)
	} else {
		user, err = s.userRepo.FindByUsername(ctx, identifier)
	}

	if err != nil {
		return "", "", "", apierror.ErrInvalidCredentials
	}

	// 2. 将提供的主密钥哈希与存储的哈希进行比较。
	// 注意：这是一个直接的字符串比较。在实际场景中，
	// 您应该使用恒定时间比较函数来防止时序攻击。
	if user.AuthHash != masterKeyHash {
		return "", "", "", apierror.ErrInvalidCredentials
	}

	// 3. 生成 JWT。
	token, err := crypto.GenerateJWT(user.ID, s.cfg.JWTSecret, s.cfg.JWTExpiration)
	if err != nil {
		return "", "", "", apierror.ErrInternalServer
	}

	// 4. 返回令牌和主盐。
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
	fmt.Printf("Verification code for %s: %s\n", emailAddr, code)

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
			// 在生产环境中，这里应该有更健壮的日志记录
			fmt.Printf("Failed to send verification email to %s: %v\n", emailAddr, err)
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
			fmt.Printf("Password reset requested for non-existent user: %s\n", emailAddr)
			return nil
		}
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
		fmt.Printf("ResetPassword link for %s: %s\n", emailAddr, resetLink)

		err := s.emailSvc.SendPasswordResetEmail(user.Email, resetLink)
		if err != nil {
			fmt.Printf("Failed to send password reset email to %s: %v\n", user.Email, err)
		}
	}()

	return nil
}

// ResetPassword 使用有效的重置令牌重置用户的密码。
func (s *AuthService) ResetPassword(ctx context.Context, token, newMasterKeyHash, newMasterSalt string) error {
	// 1. 验证令牌。
	if token == "" {
		return apierror.ErrInvalidResetToken
	}
	tokenHash := crypto.HashString(token)

	user, err := s.userRepo.FindByResetPasswordToken(ctx, tokenHash)
	if err != nil {
		if err == core.ErrUserNotFound {
			return apierror.ErrInvalidResetToken
		}
		return apierror.ErrInternalServer
	}

	// 2. 检查令牌是否已过期。
	if user.ResetPasswordTokenExpiresAt == nil || time.Now().After(*user.ResetPasswordTokenExpiresAt) {
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
		return apierror.ErrInternalServer
	}

	return nil
}
