package email

import (
	"bytes"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

// EmailService 定义了发送邮件的接口。
type EmailService interface {
	SendEmail(to, subject, body string) error
	SendPasswordResetEmail(to, resetLink string) error
	SendVerificationCodeEmail(to, code string) error
}

// SMTPEmailService 是 EmailService 的一个实现，使用 SMTP 发送邮件。
type SMTPEmailService struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// NewSMTPEmailService 创建一个新的 SMTPEmailService。
func NewSMTPEmailService(host string, port int, username, password, from string) *SMTPEmailService {
	return &SMTPEmailService{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

// SendEmail 发送一封纯文本邮件。
func (s *SMTPEmailService) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// PasswordResetTemplateData 是密码重置邮件模板所需的数据。
type PasswordResetTemplateData struct {
	ResetLink string
}

// SendPasswordResetEmail 发送一封密码重置邮件。
func (s *SMTPEmailService) SendPasswordResetEmail(to, resetLink string) error {
	// 在实际应用中，您会希望使用一个更美观的HTML模板。
	// 为简单起见，我们在这里使用一个基本的字符串模板。
	const templateStr = `
	<html>
	<body>
	<p>您好,</p>
	<p>我们收到了您的密码重置请求。请点击下面的链接来重置您的密码：</p>
	<p><a href="{{.ResetLink}}">重置密码</a></p>
	<p>如果您没有请求重置密码，请忽略此邮件。</p>
	<p>此链接将在1小时后失效。</p>
	<p>谢谢,</p>
	<p>EasyPassword 团队</p>
	</body>
	</html>
	`

	tmpl, err := template.New("passwordReset").Parse(templateStr)
	if err != nil {
		return fmt.Errorf("无法解析密码重置邮件模板: %w", err)
	}

	var body bytes.Buffer
	data := PasswordResetTemplateData{ResetLink: resetLink}
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("无法执行密码重置邮件模板: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "重置您的 EasyPassword 密码")
	m.SetBody("text/html", body.String())

	// 在生产环境中，您可能需要配置TLS/SSL。
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送密码重置邮件失败: %w", err)
	}

	return nil
}

// VerificationCodeTemplateData 是验证码邮件模板所需的数据。
type VerificationCodeTemplateData struct {
	Code string
}

// SendVerificationCodeEmail 发送一封包含验证码的邮件。
func (s *SMTPEmailService) SendVerificationCodeEmail(to, code string) error {
	const templateStr = `
	<html>
	<body>
	<p>您好,</p>
	<p>您的注册验证码是：</p>
	<p style="font-size: 24px; font-weight: bold; color: #333;">{{.Code}}</p>
	<p>此验证码将在5分钟后失效。</p>
	<p>如果您没有请求此验证码，请忽略此邮件。</p>
	<p>谢谢,</p>
	<p>EasyPassword 团队</p>
	</body>
	</html>
	`

	tmpl, err := template.New("verificationCode").Parse(templateStr)
	if err != nil {
		return fmt.Errorf("无法解析验证码邮件模板: %w", err)
	}

	var body bytes.Buffer
	data := VerificationCodeTemplateData{Code: code}
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("无法执行验证码邮件模板: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "您的 EasyPassword 注册验证码")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送验证码邮件失败: %w", err)
	}

	return nil
}