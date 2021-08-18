package config

const (
	RoleAdmin  = "ADMIN"
	RoleNormal = "NORMAL"
)

func GetRolesAvailable() []string {
	return []string{RoleAdmin, RoleNormal}
}
