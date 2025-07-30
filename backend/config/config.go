package config

import (
	"os"
	"strconv"
	"time"
)

// Config 保存应用程序配置。
type Config struct {
	DatabaseURL    string
	JWTSecret      string
	JWTExpiration  time.Duration
	DBType         string
	DBPath         string
	SMTPHost       string
	SMTPPort       int
	SMTPUser       string
	SMTPPassword   string
	SMTPFrom       string
	FrontendURL    string
	LogLevel       string
	LogFormat      string
}

// Load 从环境变量加载配置。
func Load() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// docker-compose 设置的默认连接字符串
		dbURL = "host=localhost user=postgres password=postgres dbname=easypassword port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "a-very-secret-key" // 开发环境默认值
	}

	jwtExpStr := os.Getenv("JWT_EXPIRATION_HOURS")
	jwtExpHours, err := strconv.Atoi(jwtExpStr)
	if err != nil || jwtExpHours <= 0 {
		jwtExpHours = 24 // 默认为 24 小时
	}

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "boltdb" // 默认为 boltdb
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "easypassword.db" // boltdb 的默认路径
	}

	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "localhost" // dev default
	}

	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil || smtpPort <= 0 {
		smtpPort = 1025 // dev default (e.g., MailHog)
	}

	smtpUser := os.Getenv("SMTP_USER")
	if smtpUser == "" {
		smtpUser = "" // dev default
	}

	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		smtpPassword = "" // dev default
	}

	smtpFrom := os.Getenv("SMTP_FROM")
	if smtpFrom == "" {
		smtpFrom = "no-reply@easypassword.com" // dev default
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173" // dev default
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "" {
		logFormat = "text"
	}

	return &Config{
		DatabaseURL:    dbURL,
		JWTSecret:      jwtSecret,
		JWTExpiration:  time.Hour * time.Duration(jwtExpHours),
		DBType:         dbType,
		DBPath:         dbPath,
		SMTPHost:       smtpHost,
		SMTPPort:       smtpPort,
		SMTPUser:       smtpUser,
		SMTPPassword:   smtpPassword,
		SMTPFrom:       smtpFrom,
		FrontendURL:    frontendURL,
		LogLevel:       logLevel,
		LogFormat:      logFormat,
	}
}