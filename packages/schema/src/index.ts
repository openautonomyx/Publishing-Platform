export interface Thing {
  '@context'?: 'https://schema.org';
  '@type': string;
  '@id'?: string;
  name?: string;
  description?: string;
  url?: string;
  image?: string | string[];
  sameAs?: string[];
}

export interface Person extends Thing {
  '@type': 'Person';
  givenName?: string;
  familyName?: string;
  jobTitle?: string;
  email?: string;
  affiliation?: Organization;
}

export interface Organization extends Thing {
  '@type': 'Organization' | 'NewsMediaOrganization' | 'ResearchOrganization' | 'EducationalOrganization';
  logo?: string;
}

export type CreativeWorkType =
  | 'CreativeWork'
  | 'Article'
  | 'AdvertiserContentArticle'
  | 'NewsArticle'
  | 'AnalysisNewsArticle'
  | 'AskPublicNewsArticle'
  | 'BackgroundNewsArticle'
  | 'OpinionNewsArticle'
  | 'Report'
  | 'ScholarlyArticle'
  | 'TechArticle'
  | 'BlogPosting'
  | 'SocialMediaPosting'
  | 'WebPage'
  | 'AboutPage'
  | 'CollectionPage'
  | 'FAQPage'
  | 'ProfilePage'
  | 'QAPage'
  | 'ItemPage'
  | 'Dataset'
  | 'DataCatalog'
  | 'Book'
  | 'Chapter'
  | 'HowTo'
  | 'Course'
  | 'LearningResource'
  | 'PresentationDigitalDocument'
  | 'DigitalDocument'
  | 'TextDigitalDocument'
  | 'ImageObject'
  | 'VideoObject'
  | 'AudioObject'
  | 'MediaObject'
  | 'Clip'
  | 'PodcastEpisode'
  | 'PodcastSeries'
  | 'Episode'
  | 'CreativeWorkSeries'
  | 'Periodical'
  | 'PublicationIssue'
  | 'PublicationVolume'
  | 'Review'
  | 'ClaimReview'
  | 'CriticReview'
  | 'Comment'
  | 'DiscussionForumPosting'
  | 'Recipe'
  | 'SoftwareSourceCode'
  | 'SoftwareApplication'
  | 'Map'
  | 'VisualArtwork'
  | 'Photograph';

export interface CreativeWork extends Thing {
  '@type': CreativeWorkType;
  headline?: string;
  alternativeHeadline?: string;
  abstract?: string;
  author?: Person | Organization | Array<Person | Organization>;
  editor?: Person | Organization | Array<Person | Organization>;
  contributor?: Person | Organization | Array<Person | Organization>;
  publisher?: Organization;
  dateCreated?: string;
  datePublished?: string;
  dateModified?: string;
  expires?: string;
  inLanguage?: string;
  keywords?: string[] | string;
  genre?: string | string[];
  isPartOf?: CreativeWork;
  hasPart?: CreativeWork | CreativeWork[];
  about?: Thing | Thing[];
  mentions?: Thing | Thing[];
  spatialCoverage?: Place | Place[];
  temporalCoverage?: string;
  license?: string;
  copyrightHolder?: Person | Organization;
  copyrightYear?: number;
}

export interface Place extends Thing {
  '@type': 'Place' | 'Country' | 'State' | 'City' | 'AdministrativeArea';
  address?: string;
  geo?: {
    '@type': 'GeoCoordinates';
    latitude: number;
    longitude: number;
  };
}

export interface Article extends CreativeWork {
  '@type':
    | 'Article'
    | 'AdvertiserContentArticle'
    | 'NewsArticle'
    | 'AnalysisNewsArticle'
    | 'AskPublicNewsArticle'
    | 'BackgroundNewsArticle'
    | 'OpinionNewsArticle'
    | 'ScholarlyArticle'
    | 'TechArticle'
    | 'BlogPosting';
  articleSection?: string;
  wordCount?: number;
}

export interface WebSite extends Thing {
  '@type': 'WebSite';
  publisher?: Organization;
  inLanguage?: string;
}

export interface Dataset extends CreativeWork {
  '@type': 'Dataset';
  distribution?: MediaObject | MediaObject[];
  variableMeasured?: string | string[];
}

export interface MediaObject extends CreativeWork {
  '@type': 'MediaObject' | 'ImageObject' | 'VideoObject' | 'AudioObject';
  contentUrl?: string;
  embedUrl?: string;
  encodingFormat?: string;
  duration?: string;
  width?: number;
  height?: number;
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

export function toJsonLd(schema: Thing | Thing[]): string {
  const value = Array.isArray(schema) ? schema.map(withContext) : withContext(schema);
  return JSON.stringify(value);
}
