/**
 * AI Service - Generate and enhance content with AI
 */

import { AxiosInstance } from 'axios';
import { AIGenerationRequest, AIGenerationResult } from '../types';

export class AIService {
  constructor(private apiClient: AxiosInstance) {}

  async generate(request: AIGenerationRequest): Promise<AIGenerationResult> {
    const response = await this.apiClient.post<AIGenerationResult>(
      '/ai/generate',
      request
    );
    return response.data;
  }

  async improve(content: string): Promise<{ improved: string }> {
    const response = await this.apiClient.post<{ improved: string }>(
      '/ai/improve',
      { content }
    );
    return response.data;
  }

  async summarize(content: string): Promise<{ summary: string }> {
    const response = await this.apiClient.post<{ summary: string }>(
      '/ai/summarize',
      { content }
    );
    return response.data;
  }

  async extractKeywords(content: string): Promise<{ keywords: string[] }> {
    const response = await this.apiClient.post<{ keywords: string[] }>(
      '/ai/keywords',
      { content }
    );
    return response.data;
  }
}
