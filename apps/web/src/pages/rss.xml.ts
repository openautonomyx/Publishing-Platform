import rss from '@astrojs/rss';
import { getCollection } from 'astro:content';

export async function GET(context) {
  const articles = (await getCollection('articles'))
    .filter((entry) => entry.data.published)
    .sort((a, b) => b.data.date.getTime() - a.data.date.getTime());

  return rss({
    title: 'Publishing Platform',
    description: 'AI-native professional publishing infrastructure.',
    site: context.site ?? 'https://example.com',
    items: articles.map((article) => ({
      title: article.data.title,
      description: article.data.description,
      pubDate: article.data.date,
      link: `/articles/${article.data.slug}/`
    }))
  });
}
