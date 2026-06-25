/**
 * Templates Service - Browse and use design templates
 */

import { AxiosInstance } from 'axios';
import { Template, TemplateSearchOptions, PaginatedResponse } from '../types';

export class TemplateService {
  constructor(private apiClient: AxiosInstance) {}

  async search(
    options: TemplateSearchOptions
  ): Promise<PaginatedResponse<Template>> {
    const response = await this.apiClient.get<PaginatedResponse<Template>>(
      '/templates',
      { params: options }
    );
    return response.data;
  }

  async get(id: string): Promise<Template> {
    const response = await this.apiClient.get<Template>(`/templates/${id}`);
    return response.data;
  }

  async clone(id: string): Promise<Template> {
    const response = await this.apiClient.post<Template>(
      `/templates/${id}/clone`
    );
    return response.data;
  }
}
