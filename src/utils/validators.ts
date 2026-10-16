/**
 * Input validation utilities
 */

import { ValidationError, ContentCreateRequest } from '../index';

export function validateEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

export function validateContentRequest(request: ContentCreateRequest): void {
  if (!request.title || request.title.trim() === '') {
    throw new ValidationError('Title is required and cannot be empty');
  }

  if (!request.content || request.content.trim() === '') {
    throw new ValidationError('Content is required and cannot be empty');
  }

  if (request.title.length > 500) {
    throw new ValidationError('Title must be less than 500 characters');
  }

  if (request.content.length < 10) {
    throw new ValidationError('Content must be at least 10 characters');
  }
}

export function validateLanguages(languages?: string[]): void {
  const validLanguages = [
    'en', 'es', 'fr', 'de', 'it', 'pt', 'ru', 'ja', 'zh', 'ko',
    'ar', 'hi', 'th', 'vi', 'tr', 'pl', 'nl', 'sv', 'da', 'no',
    'fi', 'id', 'ms', 'tl', 'bg', 'cs', 'hr', 'ro', 'sk', 'sl',
    'hu', 'el', 'uk', 'he', 'fa', 'ur', 'bn', 'ta', 'te', 'kn',
    'ml', 'af', 'sq', 'et', 'lt', 'lv', 'mk', 'zip'
  ];

  if (languages) {
    for (const lang of languages) {
      if (!validLanguages.includes(lang)) {
        throw new ValidationError(`Invalid language code: ${lang}`);
      }
    }
  }
}

export function validatePlatforms(platforms?: string[]): void {
  const validPlatforms = ['instagram', 'twitter', 'linkedin', 'facebook', 'tiktok', 'youtube'];

  if (platforms) {
    for (const platform of platforms) {
      if (!validPlatforms.includes(platform)) {
        throw new ValidationError(`Invalid platform: ${platform}`);
      }
    }
  }
}
