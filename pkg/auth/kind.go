package auth

// Kind represents the authentication method used in the system.
type Kind int

const (
	// HMAC represents HMAC-based authentication.
	HMAC Kind = iota + 1
	// JWT represents JSON Web Token-based authentication.
	JWT
)

// kindNames maps Kind values to their string representations.
var kindNames = map[Kind]string{
	HMAC: "hmac",
	JWT:  "jwt",
}

// String returns the string representation of the Kind.
// If the Kind is unknown, it returns "unknown".
func (k Kind) String() string {
	if name, ok := kindNames[k]; ok {
		return name
	}
	return "unknown"
}

// KindIsValid checks whether the provided Kind is a known authentication kind.
func KindIsValid(k Kind) bool {
	_, ok := kindNames[k]
	return ok
}
