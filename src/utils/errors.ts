/**
 * Error handling utilities
 */

import { AxiosError } from 'axios';
import { OpenAutonomyXError, APIError } from '../types';

export function handleAPIError(error: unknown): never {
  if (error instanceof AxiosError) {
    const status = error.response?.status || 500;
    const data = error.response?.data as any;

    const apiError: APIError = {
      code: data?.code || 'UNKNOWN_ERROR',
      message: data?.message || error.message,
      statusCode: status,
      details: data?.details
    };

    throw new OpenAutonomyXError(
      apiError.code,
      status,
      apiError.message,
      apiError.details
    );
  }

  throw new OpenAutonomyXError(
    'UNKNOWN_ERROR',
    500,
    error instanceof Error ? error.message : 'An unknown error occurred'
  );
}

export class ValidationError extends OpenAutonomyXError {
  constructor(message: string, details?: Record<string, any>) {
    super('VALIDATION_ERROR', 400, message, details);
    this.name = 'ValidationError';
  }
}

export class AuthenticationError extends OpenAutonomyXError {
  constructor(message = 'Authentication failed') {
    super('AUTHENTICATION_ERROR', 401, message);
    this.name = 'AuthenticationError';
  }
}

export class AuthorizationError extends OpenAutonomyXError {
  constructor(message = 'You do not have permission to perform this action') {
    super('AUTHORIZATION_ERROR', 403, message);
    this.name = 'AuthorizationError';
  }
}

export class NotFoundError extends OpenAutonomyXError {
  constructor(resource: string) {
    super('NOT_FOUND', 404, `${resource} not found`);
    this.name = 'NotFoundError';
  }
}

export class ConflictError extends OpenAutonomyXError {
  constructor(message: string) {
    super('CONFLICT', 409, message);
    this.name = 'ConflictError';
  }
}

export class RateLimitError extends OpenAutonomyXError {
  constructor(message = 'Rate limit exceeded') {
    super('RATE_LIMIT_ERROR', 429, message);
    this.name = 'RateLimitError';
  }
}
