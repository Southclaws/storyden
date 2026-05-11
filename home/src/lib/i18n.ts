import { defineI18n } from "fumadocs-core/i18n";

export const i18n = defineI18n({
  defaultLanguage: "en",
  languages: ["en", "zh"],
  hideLocale: "default-locale",
});

export type Locale = "en" | "zh";

export function getLocaleFromSlug(slug?: string[]): Locale {
  return slug?.[0] === "zh" ? "zh" : "en";
}

export function isChineseLocale(locale: Locale) {
  return locale === "zh";
}
