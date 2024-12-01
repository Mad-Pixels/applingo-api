package auth

type Role int

const (
	Guest Role = iota + 1
	Device
	User
	SuperUser
	Manager
	Admin
	SuperAdmin
)

var RoleNames = map[Role]string{
	Guest:      "guest",
	Device:     "device",
	User:       "user",
	SuperUser:  "superuser",
	Manager:    "manager",
	Admin:      "admin",
	SuperAdmin: "superadmin",
}

var RolePermissions = map[Role]int{
	Guest:      1,
	Device:     3,
	User:       5,
	SuperUser:  7,
	Manager:    10,
	Admin:      15,
	SuperAdmin: 20,
}

func ParseRole(role string) (Role, bool) {
	for r, name := range RoleNames {
		if name == role {
			return r, true
		}
	}
	return Guest, false
}

func GetPermissionLevel(role Role) int {
	if level, exists := RolePermissions[role]; exists {
		return level
	}
	return RolePermissions[Guest]
}
