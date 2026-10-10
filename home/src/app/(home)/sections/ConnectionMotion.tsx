"use client";

import { type ReactNode, useEffect, useRef } from "react";

/** Pause decorative flow offscreen, in hidden tabs, and for reduced motion. */
export function ConnectionMotion({ children }: { children: ReactNode }) {
  const root = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const element = root.current;
    if (!element) return;
    const preference = matchMedia("(prefers-reduced-motion: reduce)");
    let inView = false;
    const updatePlayback = () => {
      element.dataset.flowActive = String(
        inView && !preference.matches && !document.hidden,
      );
    };
    preference.addEventListener("change", updatePlayback);
    document.addEventListener("visibilitychange", updatePlayback);
    const observer = new IntersectionObserver(
      ([entry]) => {
        inView = !!entry?.isIntersecting && entry.intersectionRatio >= 0.6;
        updatePlayback();
      },
      { threshold: 0.6 },
    );
    observer.observe(element);
    return () => {
      observer.disconnect();
      preference.removeEventListener("change", updatePlayback);
      document.removeEventListener("visibilitychange", updatePlayback);
      element.dataset.flowActive = "false";
    };
  }, []);

  return (
    <div className="connection-motion" ref={root}>
      {children}
    </div>
  );
}
