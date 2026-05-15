import { getDb } from './client';
import type { Role } from '@publishing-platform/auth';

export async function createMembership(input: {
  tenantId: string;
  userId: string;
  role: Role;
}) {
  const db = await getDb();
  return db.create('membership', {
    tenant: `tenant:${input.tenantId}`,
    user: `user:${input.userId}`,
    role: input.role,
    created_at: new Date().toISOString(),
  });
}

export async function getMembership(tenantId: string, userId: string) {
  const db = await getDb();
  const [rows] = await db.query<[unknown[]]>(
    `SELECT * FROM membership
     WHERE tenant = type::thing('tenant', $tenantId)
       AND user = type::thing('user', $userId)
     LIMIT 1;`,
    { tenantId, userId },
  );

  return rows?.[0] ?? null;
}
