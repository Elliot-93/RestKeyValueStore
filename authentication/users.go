package authentication

import (
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	Admin              = "admin"
	AuthUsernameCtxKey = "AuthenticatedUsernameContextKey"
	CookieName         = "Authorization"
	HeaderName         = "Authorization"
)

var JwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var usersAndPasswords = map[string]string{
	"user_a": "$2a$14$VbJa7E/fkF3Tv4Xq3h1YvO7JlMg3B4FMd1fLoKBXXNU.0wc5HDEay",
	"user_b": "$2a$14$9kLyoUWv6MLakeYBnM3JKOLqDDQd7H.S9jWczLyDdvYwFfPmiF4N2",
	"user_c": "$2a$14$unbsetdEc3S9fzRk11.L..QHjMwcajQ7HfyZBr6LL.uxA/rfTmnNK",
	"admin":  "$2a$14$dfLFyLby9AXiRsF7f6NEyuZNMtv6WStPQwdX0gcJGK9dscMAiEorG",
}

func CheckPassword(username, password string) bool {
	hash := usersAndPasswords[username]
	if hash == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
