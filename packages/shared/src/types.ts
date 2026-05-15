export type Role = 'owner' | 'admin' | 'editor' | 'author' | 'viewer';

export type ArticleStatus =
  | 'draft'
  | 'in_review'
  | 'approved'
  | 'scheduled'
  | 'published'
  | 'archived';

export interface Tenant {
  id: string;
  name: string;
  slug: string;
}

export interface Article {
  id: string;
  tenant: string;
  title: string;
  slug: string;
  status: ArticleStatus;
  content: Record<string, unknown>;
  publishedAt?: string;
  updatedAt: string;
}
