export const I18N_COOKIE_NAME = "storyden_locale";

export const locales = ["en", "zh"] as const;

export type Locale = (typeof locales)[number];

export const defaultLocale: Locale = "en";

export function isLocale(value: string | null | undefined): value is Locale {
  return value === "en" || value === "zh";
}

export function normalizeLocale(value: string | null | undefined): Locale {
  if (isLocale(value)) {
    return value;
  }

  return defaultLocale;
}
