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