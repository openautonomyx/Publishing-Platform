/**
 * Multilingual Service - Generate and translate content in 50+ languages
 */

import { AxiosInstance } from 'axios';
import {
  MultilingualContent,
  TranslationRequest,
  TranslationResult,
  SupportedLanguage
} from '../types';

export class MultilingualService {
  constructor(private apiClient: AxiosInstance) {}

  async generateMultilingual(
    topic: string,
    languages: SupportedLanguage[]
  ): Promise<MultilingualContent[]> {
    const response = await this.apiClient.post<MultilingualContent[]>(
      '/multilingual/generate',
      { topic, languages }
    );
    return response.data;
  }

  async translate(request: TranslationRequest): Promise<TranslationResult> {
    const response = await this.apiClient.post<TranslationResult>(
      '/multilingual/translate',
      request
    );
    return response.data;
  }

  async detectLanguage(content: string): Promise<{ language: SupportedLanguage }> {
    const response = await this.apiClient.post<{ language: SupportedLanguage }>(
      '/multilingual/detect',
      { content }
    );
    return response.data;
  }

  async getAvailableLanguages(): Promise<{ languages: SupportedLanguage[] }> {
    const response = await this.apiClient.get<{ languages: SupportedLanguage[] }>(
      '/multilingual/languages'
    );
    return response.data;
  }
}
