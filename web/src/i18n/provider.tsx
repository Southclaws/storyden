"use client";

import {
  PropsWithChildren,
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

import {
  I18N_COOKIE_NAME,
  Locale,
  defaultLocale,
  isLocale,
  normalizeLocale,
} from "./config";
import { Translate, TranslationParams, interpolate } from "./format";
import { messages } from "./resources";

type I18nValue = {
  locale: Locale;
  setLocale: (next: Locale) => void;
  t: Translate;
};

const I18nContext = createContext<I18nValue | null>(null);

type Props = PropsWithChildren<{
  initialLocale: Locale;
}>;

export function I18nProvider({ initialLocale, children }: Props) {
  const [locale, setLocaleState] = useState<Locale>(
    normalizeLocale(initialLocale),
  );

  useEffect(() => {
    const normalized = normalizeLocale(initialLocale);
    setLocaleState((previous) =>
      previous === normalized ? previous : normalized,
    );
  }, [initialLocale]);

  useEffect(() => {
    const local =
      typeof window !== "undefined"
        ? localStorage.getItem(I18N_COOKIE_NAME)
        : null;

    if (isLocale(local) && local !== locale) {
      setLocaleState(local);
    }
    // Only hydrate from client storage once after mount to avoid update loops.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    document.documentElement.lang = locale;
    localStorage.setItem(I18N_COOKIE_NAME, locale);
    document.cookie = `${I18N_COOKIE_NAME}=${locale}; path=/; max-age=31536000; samesite=lax`;
  }, [locale]);

  const value = useMemo<I18nValue>(() => {
    const t: Translate = (key, params) => {
      const dictionary = messages[locale] ?? messages[defaultLocale];
      const message = dictionary[key] ?? messages[defaultLocale][key] ?? key;
      return interpolate(message, params);
    };

    return {
      locale,
      setLocale: (next: Locale) =>
        setLocaleState((previous) => (previous === next ? previous : next)),
      t,
    };
  }, [locale]);

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
}

export function useI18n() {
  const value = useContext(I18nContext);

  if (!value) {
    return {
      locale: defaultLocale,
      setLocale: () => undefined,
      t: (key: string, params?: TranslationParams) =>
        interpolate(messages[defaultLocale][key] ?? key, params),
    };
  }

  return value;
}
