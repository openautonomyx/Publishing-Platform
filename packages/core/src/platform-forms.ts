export const PLATFORM_FORMS = ['ui', 'openapi', 'webhooks', 'oci'] as const;

export type PlatformForm = (typeof PLATFORM_FORMS)[number];

export function listPlatformForms(): PlatformForm[] {
  return [...PLATFORM_FORMS];
}

export function hasPlatformForm(form: string): form is PlatformForm {
  return PLATFORM_FORMS.includes(form as PlatformForm);
}
