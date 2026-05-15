import { getDb } from './client';

export async function createTenant(input: {
  name: string;
  slug: string;
  domain?: string;
  settings?: Record<string, unknown>;
}) {
  const db = await getDb();
  return db.create('tenant', {
    name: input.name,
    slug: input.slug,
    domain: input.domain,
    settings: input.settings ?? {},
    created_at: new Date().toISOString(),
  });
}

export async function getTenantBySlug(slug: string) {
  const db = await getDb();
  const [tenant] = await db.query<[unknown[]]>(
    'SELECT * FROM tenant WHERE slug = $slug LIMIT 1;',
    { slug },
  );

  return tenant?.[0] ?? null;
}
