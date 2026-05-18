import { defineCollection, z } from 'astro:content';

const articles = defineCollection({
  type: 'content',
  schema: z.object({
    title: z.string(),
    slug: z.string(),
    description: z.string(),
    author: z.string(),
    date: z.coerce.date(),
    tags: z.array(z.string()).default([]),
    categories: z.array(z.string()).default([]),
    published: z.boolean().default(false),
    schemaType: z.string().default('BlogPosting'),
    coverImage: z.string().optional()
  })
});

export const collections = { articles };
