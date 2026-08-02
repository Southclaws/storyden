"use client";

import { useRouter } from "next/navigation";
import { useCallback } from "react";

import { WEB_ADDRESS } from "@/config";

type NavigationHistory = {
  currentEntry: { index: number } | null;
  entries(): readonly { url: string | null }[];
};

function getPreviousEntry() {
  const navigation = (
    window as typeof window & { navigation?: NavigationHistory }
  ).navigation;
  const currentIndex = navigation?.currentEntry?.index;

  if (
    navigation === undefined ||
    currentIndex === undefined ||
    currentIndex < 1
  ) {
    return undefined;
  }

  return navigation.entries()[currentIndex - 1];
}

function isStorydenURL(value: string) {
  try {
    return new URL(value).origin === new URL(WEB_ADDRESS).origin;
  } catch {
    return false;
  }
}

export function useSmartBack(fallbackHref = "/") {
  const router = useRouter();

  return useCallback(() => {
    const previous = getPreviousEntry();

    if (previous?.url && isStorydenURL(previous.url)) {
      router.back();
      return;
    }

    router.push(fallbackHref);
  }, [fallbackHref, router]);
}
