package authentication

import "github.com/golang-jwt/jwt"

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

//todo: hash passwords
var UsersAndPasswords = map[string]string{
	"user_a": "passwordA",
	"user_b": "passwordB",
	"user_c": "passwordC",
	"admin":  "Password1",
}
