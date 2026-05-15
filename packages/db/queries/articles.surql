-- Article lifecycle queries

-- Create article
CREATE article CONTENT {
  tenant: $tenant,
  author: $author,
  title: $title,
  slug: $slug,
  status: 'draft',
  content: $content,
  language: $language ?? 'en',
  location: $location ?? NONE,
  seo: $seo ?? {},
  schema: $schema ?? {},
  updated_at: time::now()
};

-- Get public article by tenant and slug
SELECT * FROM article
WHERE tenant = $tenant
  AND slug = $slug
  AND status = 'published'
LIMIT 1;

-- Transition article status
UPDATE $article MERGE {
  status: $to_status,
  published_at: IF $to_status = 'published' THEN time::now() ELSE published_at END,
  updated_at: time::now()
};

CREATE publishing_event CONTENT {
  tenant: $tenant,
  article: $article,
  actor: $actor ?? NONE,
  from_status: $from_status ?? NONE,
  to_status: $to_status,
  created_at: time::now()
};

-- Save article revision
CREATE article_revision CONTENT {
  tenant: $tenant,
  article: $article,
  actor: $actor,
  content: $content,
  created_at: time::now()
};
