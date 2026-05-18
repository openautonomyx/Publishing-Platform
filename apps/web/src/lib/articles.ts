import type { CollectionEntry } from 'astro:content';

export type Article = CollectionEntry<'articles'>;

export function getReadingTime(body: string): string {
  const words = body.trim().split(/\s+/).filter(Boolean).length;
  const minutes = Math.max(1, Math.ceil(words / 200));
  return `${minutes} min read`;
}

export function toSlug(value: string): string {
  return value
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');
}

export function getRelatedArticles(article: Article, articles: Article[], limit = 3): Article[] {
  const tags = new Set(article.data.tags);
  const categories = new Set(article.data.categories);

  return articles
    .filter((candidate) => candidate.id !== article.id && candidate.data.published)
    .map((candidate) => {
      const tagScore = candidate.data.tags.filter((tag) => tags.has(tag)).length;
      const categoryScore = candidate.data.categories.filter((category) => categories.has(category)).length;
      return { article: candidate, score: tagScore + categoryScore };
    })
    .filter((entry) => entry.score > 0)
    .sort((a, b) => b.score - a.score || b.article.data.date.getTime() - a.article.data.date.getTime())
    .slice(0, limit)
    .map((entry) => entry.article);
}
