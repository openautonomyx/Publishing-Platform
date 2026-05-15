export interface Thing {
  '@context'?: 'https://schema.org';
  '@type': string;
  '@id'?: string;
  name?: string;
  description?: string;
  url?: string;
}

export interface Person extends Thing {
  '@type': 'Person';
  givenName?: string;
  familyName?: string;
  jobTitle?: string;
  image?: string;
  sameAs?: string[];
}

export interface Organization extends Thing {
  '@type': 'Organization';
  logo?: string;
  sameAs?: string[];
}

export interface CreativeWork extends Thing {
  headline?: string;
  author?: Person | Organization;
  publisher?: Organization;
  datePublished?: string;
  dateModified?: string;
  inLanguage?: string;
}

export interface Article extends CreativeWork {
  '@type': 'Article' | 'NewsArticle' | 'OpinionNewsArticle';
  articleSection?: string;
  keywords?: string[];
}

export interface WebSite extends Thing {
  '@type': 'WebSite';
}

export interface Dataset extends CreativeWork {
  '@type': 'Dataset';
}

export interface BreadcrumbList extends Thing {
  '@type': 'BreadcrumbList';
  itemListElement: Array<{
    '@type': 'ListItem';
    position: number;
    name: string;
    item: string;
  }>;
}

export function withContext<T extends Thing>(schema: T): T {
  return {
    '@context': 'https://schema.org',
    ...schema,
  };
}

export function toJsonLd(schema: Thing): string {
  return JSON.stringify(withContext(schema));
}
