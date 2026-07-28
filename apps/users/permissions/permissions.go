package permissions

// ثوابت الصلاحيات
const (
	// صلاحيات المستخدمين
	PermUsersView   = "users.view"
	PermUsersCreate = "users.create"
	PermUsersEdit   = "users.edit"
	PermUsersDelete = "users.delete"

	// صلاحيات الأدوار
	PermRolesView   = "roles.view"
	PermRolesCreate = "roles.create"
	PermRolesEdit   = "roles.edit"
	PermRolesDelete = "roles.delete"

	// صلاحيات الصلاحيات
	PermPermissionsView   = "permissions.view"
	PermPermissionsCreate = "permissions.create"
	PermPermissionsEdit   = "permissions.edit"
	PermPermissionsDelete = "permissions.delete"
	PermPermissionsAssign = "permissions.assign"

	// صلاحيات النظام
	PermSystemSettings = "system.settings"
	PermSystemApps     = "system.apps"
	PermSystemBackup   = "system.backup"
	PermSystemRestore  = "system.restore"
)

// AllPermissions جميع الصلاحيات
var AllPermissions = []string{
	PermUsersView,
	PermUsersCreate,
	PermUsersEdit,
	PermUsersDelete,
	PermRolesView,
	PermRolesCreate,
	PermRolesEdit,
	PermRolesDelete,
	PermPermissionsView,
	PermPermissionsCreate,
	PermPermissionsEdit,
	PermPermissionsDelete,
	PermPermissionsAssign,
	PermSystemSettings,
	PermSystemApps,
	PermSystemBackup,
	PermSystemRestore,
}

// PermissionGroups مجموعات الصلاحيات
var PermissionGroups = map[string][]string{
	"users": {
		PermUsersView,
		PermUsersCreate,
		PermUsersEdit,
		PermUsersDelete,
	},
	"roles": {
		PermRolesView,
		PermRolesCreate,
		PermRolesEdit,
		PermRolesDelete,
	},
	"permissions": {
		PermPermissionsView,
		PermPermissionsCreate,
		PermPermissionsEdit,
		PermPermissionsDelete,
		PermPermissionsAssign,
	},
	"system": {
		PermSystemSettings,
		PermSystemApps,
		PermSystemBackup,
		PermSystemRestore,
	},
}
