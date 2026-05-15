export const PLATFORM_NATIVE_FEATURES = {
  aiNative: 'AI-assisted drafting, metadata, schema, localization, enrichment, and multimodal workflows.',
  cloudNative: 'Built for modern cloud deployment, multi-tenancy, custom domains, storage, and APIs.',
  automationNative: 'OpenAPI, MCP, triggers, webhooks, and tenant-owned automation endpoints.',
  feedNative: 'RSS, Atom, JSON Feed, and tenant/topic/author/language/location feeds.',
  canvasNative: 'Fabric.js powered visual composition for publishing assets and article blocks.',
  multimodalNative: 'Support for text, reports, PDFs, charts, images, tables, datasets, embeds, and canvas graphics.',
  multilingualNative: 'Language-specific versions, locale routing, translation workflows, and hreflang.',
  locationPrimitive: 'Plus Codes, Schema.org Place, geo coordinates, and linked SurrealDB geo objects.',
  timePrimitive: 'Published time, updated time, scheduled time, event time, report periods, and archives.',
  schemaOrgNative: 'Schema.org entities and JSON-LD for public semantic representation.',
  seoNative: 'Metadata, canonical URLs, sitemaps, robots, hreflang, feed discovery, and structured data.',
  composableNative: 'Modular packages, APIs, adapters, and extensible platform primitives.',
  professionalPublishingNative: 'Authoritative publishing workflows for experts, institutions, and research teams.',
  astroNative: 'Astro as the lightweight publishing frontend with React islands for interactive tools.',
  surrealDbNative: 'SurrealDB-backed tenancy, publishing flow, sessions, RBAC, schema entities, and geo objects.',
  authNative: 'Auth.js / NextAuth-compatible authentication with SurrealDB-backed app identity.',
  rbacNative: 'Owner, admin, editor, author, and viewer roles with explicit permissions.',
  apiNative: 'OpenAPI, MCP, webhooks, triggers, scoped API clients, and automation events.',
} as const;

export type PlatformNativeFeature = keyof typeof PLATFORM_NATIVE_FEATURES;

export function hasNativeFeature(feature: PlatformNativeFeature): boolean {
  return feature in PLATFORM_NATIVE_FEATURES;
}

export function listNativeFeatures(): Array<{
  key: PlatformNativeFeature;
  description: (typeof PLATFORM_NATIVE_FEATURES)[PlatformNativeFeature];
}> {
  return Object.entries(PLATFORM_NATIVE_FEATURES).map(([key, description]) => ({
    key: key as PlatformNativeFeature,
    description,
  }));
}
