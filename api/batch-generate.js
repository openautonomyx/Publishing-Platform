import { aiJsonToMarkdown, renderArticleHtml, slugify } from '../lib/publishing-utils.js';

export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  const { items = [] } = req.body || {};

  const results = items.map((item, index) => {
    const title = item.title || `Generated Article ${index + 1}`;
    const suggestedSlug = item.suggestedSlug || slugify(title);
    const meta = {
      ...item,
      title,
      suggestedSlug,
      canonicalUrl: item.canonicalUrl || `https://example.com/posts/${suggestedSlug}/`
    };

    return {
      title,
      slug: suggestedSlug,
      markdown: aiJsonToMarkdown(meta),
      html: renderArticleHtml(meta)
    };
  });

  res.status(200).json({ count: results.length, results });
}
