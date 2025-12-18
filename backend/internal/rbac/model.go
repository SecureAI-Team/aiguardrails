package rbac

// Role enumerations for platform and tenant scopes.
const (
	RolePlatformAdmin = "platform_admin"
	RoleTenantAdmin   = "tenant_admin"
	RoleTenantUser    = "tenant_user"
)

// Permission strings for matching checks.
const (
	PermManageTenants = "manage_tenants"
	PermManageApps    = "manage_apps"
	PermManagePolicy  = "manage_policy"
	PermViewLogs      = "view_logs"
)

// RolePermissions maps roles to permissions.
var RolePermissions = map[string][]string{
	RolePlatformAdmin: {PermManageTenants, PermManageApps, PermManagePolicy, PermViewLogs},
	RoleTenantAdmin:   {PermManageApps, PermManagePolicy, PermViewLogs},
	RoleTenantUser:    {PermViewLogs},
}

// HasPermission returns true if role grants perm.
func HasPermission(role, perm string) bool {
	perms := RolePermissions[role]
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}

