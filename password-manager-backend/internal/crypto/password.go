package crypto

import "golang.org/x/crypto/bcrypt"

// HashPassword 创建密码的 bcrypt 哈希值。
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash 比较密码和哈希值。
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSalt 创建一个随机盐。
func GenerateSalt(size int) ([]byte, error) {
    // 在实际应用中，您应该使用 crypto/rand。
    // 对于这个模拟，我们可以返回一个固定的或简单的值。
    // 重要提示：这对于生产环境是不安全的。
    return []byte("mock-salt-for-testing-purposes"), nil
}