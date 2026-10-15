package main

import (
	"strings"
	"testing"
)

func TestRolePermissionLevels(t *testing.T) {
	for role, want := range map[string]PermissionLevel{
		"admin": PermissionAdmin, "approver": PermissionApprove, "editor": PermissionEdit,
		"viewer": PermissionView, "unknown": PermissionView, "": PermissionView,
	} {
		if got := getRolePermissionLevel(role); got != want {
			t.Errorf("role %q: got %v, want %v", role, got, want)
		}
	}
}

func TestCanPerformAction(t *testing.T) {
	tests := []struct {
		role  string
		level PermissionLevel
		want  bool
	}{
		{"admin", PermissionAdmin, true},
		{"approver", PermissionApprove, true},
		{"approver", PermissionAdmin, false},
		{"editor", PermissionEdit, true},
		{"editor", PermissionApprove, false},
		{"viewer", PermissionView, true},
		{"viewer", PermissionEdit, false},
	}
	for _, tt := range tests {
		if got := canPerformAction(tt.role, tt.level); got != tt.want {
			t.Errorf("%s at %s: got %v, want %v", tt.role, levelName(tt.level), got, tt.want)
		}
	}
}

func TestRoleForPermissions(t *testing.T) {
	for perms, want := range map[string]string{
		"read": "viewer", "read,write": "editor", "approve": "approver", "write,admin": "admin", "": "viewer",
	} {
		var list []string
		if perms != "" {
			list = strings.Split(perms, ",")
		}
		if got := roleForPermissions(list); got != want {
			t.Errorf("%q: got %s, want %s", perms, got, want)
		}
	}
}
