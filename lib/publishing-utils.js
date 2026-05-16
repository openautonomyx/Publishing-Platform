// Shared publishing helpers for AI-generated content.

export function slugify(value = '') {
  return String(value)
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)/g, '');
}

export function escapeHtml(value = '') {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}

export function generateSeoHead(meta = {}) {
  const title = escapeHtml(meta.seoTitle || meta.title || 'Publishing Platform');
  const description = escapeHtml(meta.metaDescription || meta.description || 'AI-native publishing content.');
  const canonicalUrl = escapeHtml(meta.canonicalUrl || 'https://openautonomyx.github.io/Publishing-Platform/');
  const ogTitle = escapeHtml(meta.ogTitle || title);
  const ogDescription = escapeHtml(meta.ogDescription || description);
  const ogImage = escapeHtml(meta.ogImage || `${canonicalUrl.replace(/\/$/, '')}/og-image.png`);
  const ogImageAlt = escapeHtml(meta.ogImageAlt || `${title} preview image`);
  const tags = Array.isArray(meta.tags) ? meta.tags.join(', ') : meta.tags || '';

  return `<title>${title}</title>
<meta name="description" content="${description}" />
<meta name="keywords" content="${escapeHtml(tags)}" />
<link rel="canonical" href="${canonicalUrl}" />
<meta property="og:title" content="${ogTitle}" />
<meta property="og:description" content="${ogDescription}" />
<meta property="og:type" content="article" />
<meta property="og:url" content="${canonicalUrl}" />
<meta property="og:image" content="${ogImage}" />
<meta property="og:image:alt" content="${ogImageAlt}" />
<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:title" content="${ogTitle}" />
<meta name="twitter:description" content="${ogDescription}" />
<meta name="twitter:image" content="${ogImage}" />`;
}

export function generateSchema(meta = {}) {
  return {
    '@context': 'https://schema.org',
    '@type': 'Article',
    headline: meta.title || meta.seoTitle || 'Publishing Platform Article',
    description: meta.metaDescription || meta.description || '',
    author: {
      '@type': 'Organization',
      name: meta.author || 'OpenAutonomyX'
    },
    publisher: {
      '@type': 'Organization',
      name: 'OpenAutonomyX'
    },
    mainEntityOfPage: meta.canonicalUrl || '',
    url: meta.canonicalUrl || '',
    image: meta.ogImage || undefined,
    keywords: Array.isArray(meta.tags) ? meta.tags : [],
    articleSection: meta.category || 'Publishing',
    datePublished: meta.datePublished || new Date().toISOString().slice(0, 10),
    dateModified: meta.dateModified || new Date().toISOString().slice(0, 10)
  };
}

export function generateSchemaScript(meta = {}) {
  return `<script type="application/ld+json">${JSON.stringify(generateSchema(meta), null, 2)}</script>`;
}

export function aiJsonToMarkdown(meta = {}) {
  const title = meta.title || meta.seoTitle || 'Untitled Post';
  const slug = meta.suggestedSlug || slugify(title);
  const tags = Array.isArray(meta.tags) ? meta.tags : [];
  return `---
title: "${String(title).replaceAll('"', '\\"')}"
description: "${String(meta.metaDescription || meta.description || '').replaceAll('"', '\\"')}"
date: "${meta.datePublished || new Date().toISOString().slice(0, 10)}"
author: "${meta.author || 'OpenAutonomyX'}"
category: "${meta.category || 'Publishing'}"
tags: ${JSON.stringify(tags)}
slug: "${slug}"
canonicalUrl: "${meta.canonicalUrl || ''}"
ogTitle: "${String(meta.ogTitle || title).replaceAll('"', '\\"')}"
ogDescription: "${String(meta.ogDescription || meta.metaDescription || '').replaceAll('"', '\\"')}"
ogImage: "${meta.ogImage || ''}"
ogImageAlt: "${String(meta.ogImageAlt || '').replaceAll('"', '\\"')}"
---

# ${title}

${meta.subtitle || meta.metaDescription || ''}

## Summary

${meta.designNotes || meta.imageAnalysis || 'Write your article body here.'}
`;
}

export function renderArticleHtml(meta = {}, bodyHtml = '') {
  const head = generateSeoHead(meta);
  const schema = generateSchemaScript(meta);
  const title = escapeHtml(meta.title || meta.seoTitle || 'Untitled Post');
  const category = escapeHtml(meta.category || 'Publishing');
  const date = escapeHtml(meta.datePublished || new Date().toISOString().slice(0, 10));
  const body = bodyHtml || `<p>${escapeHtml(meta.subtitle || meta.metaDescription || '')}</p>`;

  return `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  ${head}
  ${schema}
  <style>
    body{margin:0;font-family:system-ui,-apple-system,"Segoe UI",sans-serif;background:#fff;color:#0f172a}.container{width:min(760px,calc(100% - 40px));margin:0 auto;padding:40px 0 80px}a{color:#2563eb;text-decoration:none;font-weight:700}h1{font-size:clamp(2.4rem,7vw,4.5rem);line-height:.95;letter-spacing:-.06em;margin:16px 0}.meta,.back{color:#64748b}.content p{font-size:1.15rem;line-height:1.9;color:#334155}.tags{display:flex;gap:8px;flex-wrap:wrap;margin-top:28px}.tag{background:#eef2ff;color:#2563eb;border-radius:999px;padding:6px 10px;font-size:.85rem;font-weight:700}
  </style>
</head>
<body>
  <main class="container">
    <a class="back" href="index.html">← Back to Posts</a>
    <div class="meta">${date} · ${category}</div>
    <h1>${title}</h1>
    <div class="content">${body}</div>
    <div class="tags">${(meta.tags || []).map(tag => `<span class="tag">${escapeHtml(tag)}</span>`).join('')}</div>
  </main>
</body>
</html>`;
}
