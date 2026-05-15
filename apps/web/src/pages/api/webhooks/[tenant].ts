import type { APIRoute } from 'astro';

export const POST: APIRoute = async ({ params, request }) => {
  const tenant = params.tenant;
  const payload = await request.json().catch(() => null);

  return new Response(
    JSON.stringify({
      ok: true,
      tenant,
      received: payload,
    }),
    {
      status: 200,
      headers: { 'content-type': 'application/json' },
    },
  );
};
