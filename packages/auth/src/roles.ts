export const ROLES = ['owner', 'admin', 'editor', 'author', 'viewer'] as const;

export type Role = (typeof ROLES)[number];

export const ROLE_RANK: Record<Role, number> = {
  owner: 100,
  admin: 80,
  editor: 60,
  author: 40,
  viewer: 20,
};

export function hasRole(userRole: Role, requiredRole: Role): boolean {
  return ROLE_RANK[userRole] >= ROLE_RANK[requiredRole];
}
