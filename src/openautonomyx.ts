/**
 * OpenAutonomyX SDK - Main Client
 */

import axios, { AxiosInstance } from 'axios';
import { ContentService } from './services/content';
import { PublishingService } from './services/publishing';
import { AnalyticsService } from './services/analytics';
import { TeamService } from './services/team';
import { TemplateService } from './services/templates';
import { AIService } from './services/ai';
import { MultilingualService } from './services/multilingual';

export interface OpenAutonomyXConfig {
  apiKey: string;
  apiUrl?: string;
  timeout?: number;
}

/**
 * OpenAutonomyX - Main SDK Client
 *
 * @example
 * ```typescript
 * const oax = new OpenAutonomyX({
 *   apiKey: 'your-api-key',
 *   apiUrl: 'https://openautonomyx.com:3001/api/v1'
 * });
 *
 * // Create and publish content
 * const publication = await oax.content.create({
 *   title: 'My Article',
 *   content: 'Article content...',
 *   languages: ['en', 'es', 'fr']
 * });
 *
 * // Publish to platforms
 * await oax.publishing.publish(publication.id, {
 *   platforms: ['instagram', 'twitter', 'linkedin']
 * });
 * ```
 */
export class OpenAutonomyX {
  private apiClient: AxiosInstance;
  public apiUrl: string;
  public apiKey: string;

  // Services
  public content: ContentService;
  public publishing: PublishingService;
  public analytics: AnalyticsService;
  public team: TeamService;
  public templates: TemplateService;
  public ai: AIService;
  public multilingual: MultilingualService;

  constructor(config: OpenAutonomyXConfig) {
    this.apiKey = config.apiKey;
    this.apiUrl = config.apiUrl || 'https://openautonomyx.com:3001/api/v1';

    // Initialize HTTP client
    this.apiClient = axios.create({
      baseURL: this.apiUrl,
      timeout: config.timeout || 30000,
      headers: {
        'Authorization': `Bearer ${this.apiKey}`,
        'Content-Type': 'application/json',
        'User-Agent': 'openautonomyx-sdk/1.0.0'
      }
    });

    // Initialize services
    this.content = new ContentService(this.apiClient);
    this.publishing = new PublishingService(this.apiClient);
    this.analytics = new AnalyticsService(this.apiClient);
    this.team = new TeamService(this.apiClient);
    this.templates = new TemplateService(this.apiClient);
    this.ai = new AIService(this.apiClient);
    this.multilingual = new MultilingualService(this.apiClient);
  }

  /**
   * Test the connection to OpenAutonomyX API
   */
  async testConnection(): Promise<boolean> {
    try {
      const response = await this.apiClient.get('/health');
      return response.status === 200;
    } catch (error) {
      console.error('Connection test failed:', error);
      return false;
    }
  }

  /**
   * Get current user information
   */
  async getCurrentUser(): Promise<any> {
    try {
      const response = await this.apiClient.get('/user/profile');
      return response.data;
    } catch (error) {
      throw new Error(`Failed to get user profile: ${error}`);
    }
  }

  /**
   * Get usage statistics for current user/team
   */
  async getUsage(): Promise<any> {
    try {
      const response = await this.apiClient.get('/user/usage');
      return response.data;
    } catch (error) {
      throw new Error(`Failed to get usage: ${error}`);
    }
  }
}

/**
 * Factory function for creating OpenAutonomyX instances
 *
 * @example
 * ```typescript
 * const oax = createOpenAutonomyX({
 *   apiKey: process.env.OPENAUTONOMYX_API_KEY
 * });
 * ```
 */
export function createOpenAutonomyX(
  config: OpenAutonomyXConfig
): OpenAutonomyX {
  if (!config.apiKey) {
    throw new Error(
      'API key is required. Set OPENAUTONOMYX_API_KEY environment variable or pass apiKey in config.'
    );
  }
  return new OpenAutonomyX(config);
}
