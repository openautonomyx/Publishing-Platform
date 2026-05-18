import { getCollection } from 'astro:content';

export async function GET() {
  const articles = await getCollection('articles', ({ data }) => data.published);

  const items = articles.map((article) => ({
    title: article.data.title,
    description: article.data.description,
    slug: article.data.slug,
    author: article.data.author,
    tags: article.data.tags,
    categories: article.data.categories,
    date: article.data.date.toISOString()
  }));

  return new Response(JSON.stringify(items), {
    headers: {
      'Content-Type': 'application/json; charset=utf-8'
    }
  });
}
