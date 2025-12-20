"use client";

import { useEffect, useState } from "react";

function decodeHash() {
  const raw = window.location.hash;
  if (!raw || raw.length < 2) return "";
  try {
    return decodeURIComponent(raw.slice(1));
  } catch {
    return raw.slice(1);
  }
}

export function useFragmentScroll(elementId: string) {
  const [isTargeted, setIsTargeted] = useState(false);

  useEffect(() => {
    let clearTimer: number | undefined;
    let obs: MutationObserver | undefined;

    const run = () => {
      const hash = decodeHash();
      if (hash !== elementId) return;

      const el = document.getElementById(elementId);
      if (!el) return;

      // Get scroll offset - only apply on desktop (md breakpoint and up)
      // cannot use tokens for this, so hardcoding the value from design system.
      // spacing.20 = 5rem, convert to pixels (5rem * 16px = 80px)
      const scrollOffset = window.matchMedia("(min-width: 768px)").matches
        ? 80
        : 0;

      const elementPosition = el.getBoundingClientRect().top;
      const offsetPosition = elementPosition + window.scrollY - scrollOffset;

      window.scrollTo({
        top: offsetPosition,
        behavior: "auto",
      });

      setIsTargeted(false);

      requestAnimationFrame(() => setIsTargeted(true));

      window.clearTimeout(clearTimer);
      clearTimer = window.setTimeout(() => setIsTargeted(false), 1200);
      obs?.disconnect();
      obs = undefined;
    };

    // try now, otherwise wait until it appears (streaming/suspense bs)
    if (!document.getElementById(elementId) && decodeHash() === elementId) {
      obs = new MutationObserver(run);
      obs.observe(document.body, { childList: true, subtree: true });

      window.setTimeout(() => {
        obs?.disconnect();
        obs = undefined;
      }, 3000);
    }

    run();
    window.addEventListener("hashchange", run);

    return () => {
      window.removeEventListener("hashchange", run);
      obs?.disconnect();
      window.clearTimeout(clearTimer);
    };
  }, [elementId]);

  return isTargeted;
}
