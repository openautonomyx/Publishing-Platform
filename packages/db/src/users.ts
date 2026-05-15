import { getDb } from './client';

export async function createUser(input: {
  email: string;
  name?: string;
  externalAuthId?: string;
}) {
  const db = await getDb();
  return db.create('user', {
    email: input.email,
    name: input.name,
    external_auth_id: input.externalAuthId,
    created_at: new Date().toISOString(),
  });
}

export async function getUserByEmail(email: string) {
  const db = await getDb();
  const [rows] = await db.query<[unknown[]]>(
    'SELECT * FROM user WHERE email = $email LIMIT 1;',
    { email },
  );

  return rows?.[0] ?? null;
}
