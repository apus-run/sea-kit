package security

import "golang.org/x/crypto/bcrypt"

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(inputPass, validPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(validPass), []byte(inputPass))
	return err == nil
}
