import type { APIRoute } from 'astro';
import { listNativeFeatures } from '@publishing-platform/core';

export const GET: APIRoute = async () => {
  return new Response(JSON.stringify({
    features: listNativeFeatures(),
  }), {
    headers: {
      'content-type': 'application/json; charset=utf-8',
    },
  });
};
