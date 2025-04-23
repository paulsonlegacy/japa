// Password hashing, etc.
package utils

import "golang.org/x/crypto/bcrypt"


// Hash password using bcrypt algo
func HashPassword(password string) string {
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPwd)
}