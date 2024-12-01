package auth

type Type string

var (
	HMAC Type = "hmac"
	JWT  Type = "jwt"
)
