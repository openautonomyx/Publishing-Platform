/**
 * Publishing Service - Publish to social platforms and CMS
 */

import { AxiosInstance } from 'axios';
import { PublishRequest, PublishResult, SocialPlatform } from '../types';

export class PublishingService {
  constructor(private apiClient: AxiosInstance) {}

  async publish(publicationId: string, request: PublishRequest): Promise<PublishResult[]> {
    const response = await this.apiClient.post<PublishResult[]>(
      `/publications/${publicationId}/publish`,
      request
    );
    return response.data;
  }

  async schedule(publicationId: string, request: PublishRequest): Promise<PublishResult[]> {
    const response = await this.apiClient.post<PublishResult[]>(
      `/publications/${publicationId}/schedule`,
      request
    );
    return response.data;
  }

  async getStatus(publicationId: string): Promise<PublishResult[]> {
    const response = await this.apiClient.get<PublishResult[]>(
      `/publications/${publicationId}/status`
    );
    return response.data;
  }

  async cancel(publicationId: string, platform: SocialPlatform): Promise<void> {
    await this.apiClient.post(
      `/publications/${publicationId}/cancel`,
      { platform }
    );
  }
}
