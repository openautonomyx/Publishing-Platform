package main

// PermissionLevel orders the platform's roles: viewer < editor < approver < admin.
type PermissionLevel int

const (
	PermissionView PermissionLevel = iota
	PermissionEdit
	PermissionApprove
	PermissionAdmin
)

func getRolePermissionLevel(role string) PermissionLevel {
	switch role {
	case "admin":
		return PermissionAdmin
	case "approver":
		return PermissionApprove
	case "editor":
		return PermissionEdit
	default:
		return PermissionView
	}
}

func canPerformAction(userRole string, requiredLevel PermissionLevel) bool {
	return getRolePermissionLevel(userRole) >= requiredLevel
}

func levelName(l PermissionLevel) string {
	switch l {
	case PermissionAdmin:
		return "admin"
	case PermissionApprove:
		return "approver"
	case PermissionEdit:
		return "editor"
	default:
		return "viewer"
	}
}
