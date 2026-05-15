import type { Role } from './roles';

export interface AuthSession {
  userId: string;
  email: string;
  tenantId?: string;
  tenantSlug?: string;
  role?: Role;
  expiresAt?: string;
}

export function isAuthenticated(session: AuthSession | null | undefined): session is AuthSession {
  return !!session?.userId;
}
