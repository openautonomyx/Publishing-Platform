export const PLATFORM_TRIGGERS = [
  'tenant.created',
  'user.created',
  'membership.created',
  'article.created',
  'article.updated',
  'article.submitted_for_review',
  'article.approved',
  'article.scheduled',
  'article.published',
  'article.archived',
  'feed.generated',
  'schema.generated',
  'geo.created',
  'webhook.delivered',
  'webhook.failed',
] as const;

export type PlatformTrigger = (typeof PLATFORM_TRIGGERS)[number];

export interface PlatformEvent<TPayload = Record<string, unknown>> {
  id: string;
  tenantId?: string;
  type: PlatformTrigger;
  payload: TPayload;
  occurredAt: string;
}

export function isPlatformTrigger(value: string): value is PlatformTrigger {
  return PLATFORM_TRIGGERS.includes(value as PlatformTrigger);
}
