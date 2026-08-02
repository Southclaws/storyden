"use client";

import { usePathname, useRouter } from "next/navigation";
import { useCallback, useEffect, useRef } from "react";

import { WEB_ADDRESS } from "@/config";

const historyStorageKey = "storyden:navigation-history";
const historyLimit = 50;

function isStorydenURL(value: string) {
  try {
    return new URL(value, WEB_ADDRESS).origin === new URL(WEB_ADDRESS).origin;
  } catch {
    return false;
  }
}

function readHistory() {
  try {
    const value = JSON.parse(sessionStorage.getItem(historyStorageKey) ?? "[]");

    if (!Array.isArray(value)) {
      return [];
    }

    return value.filter(
      (entry): entry is string =>
        typeof entry === "string" && isStorydenURL(entry),
    );
  } catch {
    return [];
  }
}

function writeHistory(entries: string[]) {
  sessionStorage.setItem(
    historyStorageKey,
    JSON.stringify(entries.slice(-historyLimit)),
  );
}

export function NavigationHistoryTracker() {
  const pathname = usePathname();
  const mounted = useRef(false);

  useEffect(() => {
    const current = window.location.href;

    if (!mounted.current) {
      mounted.current = true;

      const referrer = document.referrer;
      writeHistory(
        referrer && referrer !== current && isStorydenURL(referrer)
          ? [referrer, current]
          : [current],
      );
      return;
    }

    const entries = readHistory();
    const previous = entries.at(-2);

    if (entries.at(-1) === current) {
      return;
    }

    if (previous === current) {
      entries.pop();
    } else {
      entries.push(current);
    }

    writeHistory(entries);
  }, [pathname]);

  return null;
}

export function useSmartBack(fallbackHref = "/") {
  const router = useRouter();

  return useCallback(() => {
    const entries = readHistory();
    const current = window.location.href;
    const previous =
      entries.at(-1) === current ? entries.at(-2) : entries.at(-1);

    if (
      window.history.length > 1 &&
      previous !== undefined &&
      isStorydenURL(previous)
    ) {
      router.back();
      return;
    }

    router.push(fallbackHref);
  }, [fallbackHref, router]);
}
