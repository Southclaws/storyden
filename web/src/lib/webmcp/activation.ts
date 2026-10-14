"use client";

import { useSyncExternalStore } from "react";

let enabled = false;
const listeners = new Set<() => void>();

// Keep tools available after the first chat activation, including across navigation.
export function enableWebMCPTools() {
  if (enabled) return;
  enabled = true;
  for (const listener of listeners) listener();
}

export function useWebMCPToolsEnabled() {
  return useSyncExternalStore(
    subscribe,
    () => enabled,
    () => false,
  );
}

function subscribe(listener: () => void) {
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
  };
}
