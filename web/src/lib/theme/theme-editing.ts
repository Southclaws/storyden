"use client";

import { useSyncExternalStore } from "react";

export const THEME_EDITING_STORAGE_KEY = "storyden:theme-editing";
const THEME_EDITING_EVENT = "storyden:theme-editing-change";

export function useThemeEditingEnabled() {
  return useSyncExternalStore(subscribe, readThemeEditingEnabled, () => false);
}

export function setThemeEditingEnabled(enabled: boolean) {
  if (enabled) {
    window.localStorage.setItem(THEME_EDITING_STORAGE_KEY, "true");
  } else {
    window.localStorage.removeItem(THEME_EDITING_STORAGE_KEY);
  }
  window.dispatchEvent(new Event(THEME_EDITING_EVENT));
}

export function parseThemeEditingFlag(value: string | null) {
  return value === "true";
}

function readThemeEditingEnabled() {
  return parseThemeEditingFlag(
    window.localStorage.getItem(THEME_EDITING_STORAGE_KEY),
  );
}

function subscribe(callback: () => void) {
  window.addEventListener("storage", callback);
  window.addEventListener(THEME_EDITING_EVENT, callback);
  return () => {
    window.removeEventListener("storage", callback);
    window.removeEventListener(THEME_EDITING_EVENT, callback);
  };
}
