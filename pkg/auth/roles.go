package auth

// Role represents a user role in the system.
type Role int

const (
	// Guest represents a user with the lowest permission level.
	Guest Role = iota + 1

	// Device represents an IoT or automated device.
	Device

	// User represents a standard authenticated user.
	User

	// SuperUser represents a power user with elevated access.
	SuperUser

	// Manager represents a role with access to administrative views.
	Manager

	// Admin represents a system administrator.
	Admin

	// SuperAdmin represents the highest-level system administrator.
	SuperAdmin
)

// RoleNames maps a Role value to its string representation.
var RoleNames = map[Role]string{
	Guest:      "guest",
	Device:     "device",
	User:       "user",
	SuperUser:  "superuser",
	Manager:    "manager",
	Admin:      "admin",
	SuperAdmin: "superadmin",
}

// RolePermissions maps a Role to its numeric permission level.
var RolePermissions = map[Role]int{
	Guest:      1,
	Device:     3,
	User:       5,
	SuperUser:  7,
	Manager:    10,
	Admin:      15,
	SuperAdmin: 20,
}

// ParseRole converts a string into a Role.
// It returns the corresponding Role and true if the string is valid,
// otherwise it returns Guest and false.
func ParseRole(role string) (Role, bool) {
	for r, name := range RoleNames {
		if name == role {
			return r, true
		}
	}
	return Guest, false
}

// GetPermissionLevel returns the permission level for the given Role.
// If the role is unknown, it defaults to the Guest permission level.
func GetPermissionLevel(role Role) int {
	if level, exists := RolePermissions[role]; exists {
		return level
	}
	return RolePermissions[Guest]
}
