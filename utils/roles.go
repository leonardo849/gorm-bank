package utils

type Role int

const (
	OWNER = iota
	MANAGER 
	CUSTOMER 
)


func GetRoleString(r Role) string {
	return [...]string{"OWNER", "MANAGER", "CUSTOMER"}[r]
}