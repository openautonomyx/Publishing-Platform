/**
 * Analytics Service - Track performance and engagement
 */

import { AxiosInstance } from 'axios';
import { AnalyticsMetrics, PerformanceReport } from '../types';

export class AnalyticsService {
  constructor(private apiClient: AxiosInstance) {}

  async getMetrics(publicationId: string): Promise<AnalyticsMetrics[]> {
    const response = await this.apiClient.get<AnalyticsMetrics[]>(
      `/publications/${publicationId}/analytics`
    );
    return response.data;
  }

  async getReport(dateRange: { start: string; end: string }): Promise<PerformanceReport> {
    const response = await this.apiClient.get<PerformanceReport>(
      '/analytics/report',
      { params: dateRange }
    );
    return response.data;
  }
}
