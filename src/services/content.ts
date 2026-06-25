/**
 * Content Service - Create and manage publications
 */

import { AxiosInstance } from 'axios';
import {
  Publication,
  ContentCreateRequest,
  PaginatedResponse,
  QueryOptions
} from '../types';

export class ContentService {
  constructor(private apiClient: AxiosInstance) {}

  /**
   * Create a new publication
   */
  async create(request: ContentCreateRequest): Promise<Publication> {
    const response = await this.apiClient.post<Publication>(
      '/publications',
      request
    );
    return response.data;
  }

  /**
   * Get a publication by ID
   */
  async get(id: string): Promise<Publication> {
    const response = await this.apiClient.get<Publication>(
      `/publications/${id}`
    );
    return response.data;
  }

  /**
   * List all publications
   */
  async list(
    options?: QueryOptions
  ): Promise<PaginatedResponse<Publication>> {
    const response = await this.apiClient.get<
      PaginatedResponse<Publication>
    >('/publications', {
      params: options
    });
    return response.data;
  }

  /**
   * Update a publication
   */
  async update(
    id: string,
    updates: Partial<ContentCreateRequest>
  ): Promise<Publication> {
    const response = await this.apiClient.patch<Publication>(
      `/publications/${id}`,
      updates
    );
    return response.data;
  }

  /**
   * Delete a publication
   */
  async delete(id: string): Promise<void> {
    await this.apiClient.delete(`/publications/${id}`);
  }

  /**
   * Archive a publication
   */
  async archive(id: string): Promise<Publication> {
    const response = await this.apiClient.post<Publication>(
      `/publications/${id}/archive`
    );
    return response.data;
  }

  /**
   * Restore an archived publication
   */
  async restore(id: string): Promise<Publication> {
    const response = await this.apiClient.post<Publication>(
      `/publications/${id}/restore`
    );
    return response.data;
  }

  /**
   * Duplicate a publication
   */
  async duplicate(id: string): Promise<Publication> {
    const response = await this.apiClient.post<Publication>(
      `/publications/${id}/duplicate`
    );
    return response.data;
  }
}
