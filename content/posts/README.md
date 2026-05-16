# Posts

Write new article drafts in this folder using Markdown.

## Suggested front matter

```md
---
title: "Your Post Title"
date: "2026-05-15"
author: "OpenAutonomyX"
category: "Product"
description: "A short summary for SEO and social previews."
slug: "your-post-title"
---

Your article content goes here.
```

## Publishing workflow

This repository currently publishes static HTML through GitHub Pages. To publish a Markdown draft:

1. Copy `templates/post-template.html`.
2. Save it as `posts/your-post-slug.html`.
3. Replace the title, date, category, description, and article body.
4. Add the article card to `posts/index.html`.
5. Add the article URL to `sitemap.xml`.

A future enhancement can add an automated Markdown-to-HTML build step.
