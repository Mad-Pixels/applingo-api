package auth

type Kind int

const (
	HMAC Kind = iota + 1
	JWT
)

var kindNames = map[Kind]string{
	HMAC: "hmac",
	JWT:  "jwt",
}

func (k Kind) String() string {
	if name, ok := kindNames[k]; ok {
		return name
	}
	return "unknown"
}

func KindIsValid(k Kind) bool {
	_, ok := kindNames[k]
	return ok
}
