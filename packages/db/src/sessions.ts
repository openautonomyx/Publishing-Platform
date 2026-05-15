import { getDb } from './client';

export async function createSession(input: {
  userId: string;
  tenantId?: string;
  tokenHash: string;
  expiresAt: string;
}) {
  const db = await getDb();
  return db.create('session', {
    user: `user:${input.userId}`,
    tenant: input.tenantId ? `tenant:${input.tenantId}` : undefined,
    token_hash: input.tokenHash,
    expires_at: input.expiresAt,
    created_at: new Date().toISOString(),
  });
}

export async function getSessionByTokenHash(tokenHash: string) {
  const db = await getDb();
  const [rows] = await db.query<[unknown[]]>(
    'SELECT * FROM session WHERE token_hash = $tokenHash LIMIT 1;',
    { tokenHash },
  );

  return rows?.[0] ?? null;
}

export async function deleteSession(sessionId: string) {
  const db = await getDb();
  return db.delete(`session:${sessionId}`);
}
