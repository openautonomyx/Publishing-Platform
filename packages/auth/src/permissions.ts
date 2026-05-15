import type { Role } from './roles';
import { hasRole } from './roles';

export type Permission =
  | 'tenant:read'
  | 'tenant:update'
  | 'users:invite'
  | 'users:manage'
  | 'article:create'
  | 'article:read'
  | 'article:update'
  | 'article:submit_review'
  | 'article:approve'
  | 'article:publish'
  | 'article:archive'
  | 'media:create'
  | 'media:delete'
  | 'settings:manage';

const ROLE_PERMISSIONS: Record<Role, Permission[]> = {
  owner: [
    'tenant:read',
    'tenant:update',
    'users:invite',
    'users:manage',
    'article:create',
    'article:read',
    'article:update',
    'article:submit_review',
    'article:approve',
    'article:publish',
    'article:archive',
    'media:create',
    'media:delete',
    'settings:manage',
  ],
  admin: [
    'tenant:read',
    'tenant:update',
    'users:invite',
    'users:manage',
    'article:create',
    'article:read',
    'article:update',
    'article:submit_review',
    'article:approve',
    'article:publish',
    'article:archive',
    'media:create',
    'media:delete',
    'settings:manage',
  ],
  editor: [
    'tenant:read',
    'article:create',
    'article:read',
    'article:update',
    'article:submit_review',
    'article:approve',
    'article:publish',
    'article:archive',
    'media:create',
  ],
  author: [
    'tenant:read',
    'article:create',
    'article:read',
    'article:update',
    'article:submit_review',
    'media:create',
  ],
  viewer: ['tenant:read', 'article:read'],
};

export function can(role: Role, permission: Permission): boolean {
  return ROLE_PERMISSIONS[role].includes(permission);
}

export function requireRole(role: Role, requiredRole: Role): void {
  if (!hasRole(role, requiredRole)) {
    throw new Error(`Requires role ${requiredRole}`);
  }
}

export function requirePermission(role: Role, permission: Permission): void {
  if (!can(role, permission)) {
    throw new Error(`Missing permission ${permission}`);
  }
}
