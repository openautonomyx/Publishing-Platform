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

export interface GeoCoordinates {
  '@type': 'GeoCoordinates';
  latitude: number;
  longitude: number;
  elevation?: number;
}

export interface PostalAddress {
  '@type': 'PostalAddress';
  streetAddress?: string;
  addressLocality?: string;
  addressRegion?: string;
  postalCode?: string;
  addressCountry?: string;
}

export type PlaceType =
  | 'Place'
  | 'Country'
  | 'State'
  | 'City'
  | 'AdministrativeArea'
  | 'Continent'
  | 'LandmarksOrHistoricalBuildings'
  | 'TouristAttraction';

export interface Place extends Thing {
  '@type': PlaceType;
  address?: PostalAddress | string;
  geo?: GeoCoordinates;
  containedInPlace?: Place;
  containsPlace?: Place | Place[];
}

export interface LinkedGeoObject {
  id: string;
  tenant?: string;
  schema: Place;
  parent?: string;
  ancestors?: string[];
  children?: string[];
  slug?: string;
  level?: 'continent' | 'country' | 'region' | 'state' | 'city' | 'locality' | 'venue' | 'custom';
  codes?: {
    iso2?: string;
    iso3?: string;
    regionCode?: string;
    localCode?: string;
  };
}

export interface Person extends Thing {
  '@type': 'Person';
  givenName?: string;
  familyName?: string;
  jobTitle?: string;
  email?: string;
  affiliation?: Organization;
  homeLocation?: Place;
  workLocation?: Place;
}

export interface Organization extends Thing {
  '@type': 'Organization' | 'NewsMediaOrganization' | 'ResearchOrganization' | 'EducationalOrganization';
  logo?: string;
  location?: Place;
  address?: PostalAddress | string;
  areaServed?: Place | Place[];
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
  contentLocation?: Place;
  locationCreated?: Place;
  license?: string;
  copyrightHolder?: Person | Organization;
  copyrightYear?: number;
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
