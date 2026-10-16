/**
 * OpenAutonomyX SDK - Type Definitions
 */

export type SupportedLanguage =
  | 'en' | 'es' | 'fr' | 'de' | 'it' | 'pt' | 'ru' | 'ja' | 'zh' | 'ko'
  | 'ar' | 'hi' | 'th' | 'vi' | 'tr' | 'pl' | 'nl' | 'sv' | 'da' | 'no'
  | 'fi' | 'id' | 'ms' | 'tl' | 'bg' | 'cs' | 'hr' | 'ro' | 'sk' | 'sl'
  | 'hu' | 'el' | 'uk' | 'he' | 'fa' | 'ur' | 'bn' | 'ta' | 'te' | 'kn'
  | 'ml' | 'af' | 'sq' | 'et' | 'lt' | 'lv' | 'mk' | 'zip';

export type SocialPlatform =
  | 'instagram' | 'twitter' | 'linkedin' | 'facebook' | 'tiktok' | 'youtube';

export type ContentType =
  | 'blog' | 'social' | 'email' | 'newsletter' | 'product' | 'announcement';

export type UserRole = 'owner' | 'admin' | 'editor' | 'viewer';

// ============================================================================
// Content Types
// ============================================================================

export interface ContentCreateRequest {
  title: string;
  content: string;
  contentType?: ContentType;
  languages?: SupportedLanguage[];
  tags?: string[];
  metadata?: Record<string, any>;
}

export interface Publication {
  id: string;
  title: string;
  content: string;
  contentType: ContentType;
  languages: SupportedLanguage[];
  status: 'draft' | 'published' | 'archived';
  createdAt: string;
  updatedAt: string;
  createdBy: string;
  metadata?: Record<string, any>;
  views?: number;
  engagements?: number;
}

// ============================================================================
// Publishing Types
// ============================================================================

export interface PublishRequest {
  platforms: SocialPlatform[];
  scheduledFor?: string; // ISO 8601 timestamp
  content?: string;
  metadata?: Record<string, any>;
}

export interface PublishResult {
  publicationId: string;
  platform: SocialPlatform;
  status: 'success' | 'failed' | 'scheduled';
  url?: string;
  message?: string;
  scheduledFor?: string;
}

// ============================================================================
// Analytics Types
// ============================================================================

export interface AnalyticsMetrics {
  publicationId: string;
  platform: SocialPlatform;
  language?: SupportedLanguage;
  impressions: number;
  engagement: number;
  clicks: number;
  shares: number;
  comments: number;
  likes: number;
  timestamp: string;
}

export interface PerformanceReport {
  totalImpressions: number;
  totalEngagement: number;
  topPerformers: Publication[];
  languagePerformance: Record<SupportedLanguage, AnalyticsMetrics>;
  platformPerformance: Record<SocialPlatform, AnalyticsMetrics>;
  period: {
    start: string;
    end: string;
  };
}

// ============================================================================
// Team Types
// ============================================================================

export interface TeamMember {
  id: string;
  email: string;
  name: string;
  role: UserRole;
  joinedAt: string;
  lastActive?: string;
}

export interface Team {
  id: string;
  name: string;
  owner: string;
  members: TeamMember[];
  createdAt: string;
}

export interface InviteRequest {
  email: string;
  role: UserRole;
}

// ============================================================================
// Template Types
// ============================================================================

export interface Template {
  id: string;
  name: string;
  category: 'marketing' | 'editorial' | 'social' | 'landing' | 'email';
  content: string;
  metadata: Record<string, any>;
  tags: string[];
  createdAt: string;
  updatedAt: string;
}

export interface TemplateSearchOptions {
  category?: Template['category'];
  search?: string;
  tags?: string[];
  limit?: number;
  offset?: number;
}

// ============================================================================
// AI Types
// ============================================================================

export interface AIGenerationRequest {
  topic: string;
  contentType: ContentType;
  tone?: 'professional' | 'casual' | 'humorous' | 'technical' | 'persuasive';
  language?: SupportedLanguage;
  keywords?: string[];
  targetAudience?: string;
}

export interface AIGenerationResult {
  content: string;
  title: string;
  summary: string;
  keywords: string[];
  tone: string;
  language: SupportedLanguage;
}

// ============================================================================
// Multilingual Types
// ============================================================================

export interface MultilingualContent {
  language: SupportedLanguage;
  content: string;
  title: string;
  summary: string;
  readingTime: number;
  characterCount: number;
}

export interface TranslationRequest {
  content: string;
  sourceLanguage?: SupportedLanguage;
  targetLanguages: SupportedLanguage[];
}

export interface TranslationResult {
  original: string;
  sourceLanguage: SupportedLanguage;
  translations: Record<SupportedLanguage, string>;
}

// ============================================================================
// Error Types
// ============================================================================

export interface APIError {
  code: string;
  message: string;
  statusCode: number;
  details?: Record<string, any>;
}

export class OpenAutonomyXError extends Error {
  constructor(
    public code: string,
    public statusCode: number,
    message: string,
    public details?: Record<string, any>
  ) {
    super(message);
    this.name = 'OpenAutonomyXError';
  }
}

// ============================================================================
// Query Types
// ============================================================================

export interface QueryOptions {
  limit?: number;
  offset?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
  filters?: Record<string, any>;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  limit: number;
  offset: number;
  hasMore: boolean;
}

// ============================================================================
// Response Types
// ============================================================================

export interface SuccessResponse<T> {
  success: true;
  data: T;
  message?: string;
}

export interface ErrorResponse {
  success: false;
  error: APIError;
}

export type APIResponse<T> = SuccessResponse<T> | ErrorResponse;
