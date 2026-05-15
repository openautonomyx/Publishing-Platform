import { getDb } from './client';

export type CreateArticleInput = {
  tenant: string;
  author: string;
  title: string;
  slug: string;
  content: Record<string, unknown>;
  language?: string;
  seo?: Record<string, unknown>;
  schema?: Record<string, unknown>;
  location?: string | null;
};

export async function createArticle(input: CreateArticleInput) {
  const db = await getDb();
  const [article] = await db.query<[unknown[]]>(
    `CREATE article CONTENT {
      tenant: type::thing('tenant', $tenant),
      author: type::thing('user', $author),
      title: $title,
      slug: $slug,
      status: 'draft',
      content: $content,
      language: $language,
      seo: $seo,
      schema: $schema,
      location: IF $location THEN type::thing('geo_object', $location) ELSE NONE END,
      updated_at: time::now()
    };`,
    {
      ...input,
      language: input.language ?? 'en',
      seo: input.seo ?? {},
      schema: input.schema ?? {},
      location: input.location ?? null,
    },
  );

  return article?.[0] ?? null;
}

export async function getPublishedArticleBySlug(tenantSlug: string, articleSlug: string) {
  const db = await getDb();
  const [rows] = await db.query<[unknown[]]>(
    `LET $tenant_record = (SELECT VALUE id FROM tenant WHERE slug = $tenantSlug LIMIT 1)[0];
     SELECT * FROM article
     WHERE tenant = $tenant_record
       AND slug = $articleSlug
       AND status = 'published'
     LIMIT 1;`,
    { tenantSlug, articleSlug },
  );

  return rows?.[0] ?? null;
}

export async function publishArticle(articleId: string, actorId: string) {
  const db = await getDb();
  const [updated] = await db.query<[unknown[]]>(
    `LET $article_record = type::thing('article', $articleId);
     LET $current = SELECT * FROM $article_record;
     UPDATE $article_record MERGE {
       status: 'published',
       published_at: time::now(),
       updated_at: time::now()
     };
     CREATE publishing_event CONTENT {
       tenant: $current[0].tenant,
       article: $article_record,
       actor: type::thing('user', $actorId),
       from_status: $current[0].status,
       to_status: 'published',
       created_at: time::now()
     };`,
    { articleId, actorId },
  );

  return updated?.[0] ?? null;
}
