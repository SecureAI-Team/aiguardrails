package auth

// MapRole maps IdP role strings to internal roles.
func MapRole(sourceRole, adminRole, userRole string) string {
	switch sourceRole {
	case adminRole:
		return "tenant_admin"
	case userRole:
		return "tenant_user"
	default:
		return sourceRole
	}
}

