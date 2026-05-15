import type { APIRoute } from 'astro';
import { listNativeFeatures } from '@publishing-platform/core';

export const GET: APIRoute = async () => {
  const features = listNativeFeatures();

  const document = {
    openapi: '3.1.0',
    info: {
      title: 'Publishing Platform API',
      version: '0.0.1',
      description: 'OpenAPI surface for the AI-native professional publishing platform.',
    },
    paths: {
      '/features.json': {
        get: {
          summary: 'List platform-native features',
          operationId: 'listNativeFeatures',
          responses: {
            '200': {
              description: 'Platform-native feature registry',
              content: {
                'application/json': {
                  schema: {
                    type: 'object',
                    properties: {
                      features: {
                        type: 'array',
                        items: { $ref: '#/components/schemas/NativeFeature' },
                      },
                    },
                    required: ['features'],
                  },
                  examples: {
                    default: {
                      value: { features },
                    },
                  },
                },
              },
            },
          },
        },
      },
    },
    components: {
      schemas: {
        NativeFeature: {
          type: 'object',
          properties: {
            key: { type: 'string' },
            description: { type: 'string' },
          },
          required: ['key', 'description'],
        },
      },
    },
  };

  return new Response(JSON.stringify(document, null, 2), {
    headers: {
      'content-type': 'application/json; charset=utf-8',
    },
  });
};
