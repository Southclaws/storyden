"use client";

import { Cache, State } from "swr";

const CACHE_KEY_NAME = "swr-cache";

export function cacheProvider<Data = any>(): () => Cache<Data> {
  const ls = localStorage.getItem(CACHE_KEY_NAME);
  const parsedCache = ls && JSON.parse(ls);

  // The cache is an in-memory map for performance.
  const map = new Map<string, State<Data>>(parsedCache);

  // When a session ends, the local cache is dumped to local storage. Note that
  // this is not compatible with multiple browser sessions at the same time, the
  // most recently closed tab will set the current cache state.
  window.addEventListener("beforeunload", () => {
    const appCache = JSON.stringify(Array.from(map.entries()));
    localStorage.setItem(CACHE_KEY_NAME, appCache);
  });

  function get(key: string): State<Data> | undefined {
    return map.get(key) as State<Data>;
  }

  function set(key: string, value: State<Data>) {
    map.set(key, value);
  }

  function deleteKey(key: string) {
    map.delete(key);
  }

  function keys(): IterableIterator<string> {
    return map.keys();
  }

  return () => ({
    get,
    set,
    ["delete"]: deleteKey,
    keys,
  });
}
