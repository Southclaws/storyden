"use client";

import { useEffect, useRef } from "react";
import { Cache, State } from "swr";

const CACHE_KEY_NAME = "swr-cache";

type CacheProvider = () => Cache;

type Data = State<any, any>;

export function useCacheProvider(): CacheProvider {
  const cache = useRef<Map<string, Data>>(new Map());

  useEffect(() => {
    const appCache = loadCache();
    if (appCache) {
      const map = new Map(JSON.parse(appCache) as Iterable<[string, Data]>);
      console.log("Loading cache", { entries: map.size });
      map.forEach((value, key) => cache.current.set(key, value));
    }

    const saveCache = () => {
      const entries = cache.current.entries();
      const a = Array.from(entries);
      console.log("Saving cache", { entries: a.length });
      storeCache(JSON.stringify(a));
    };

    window.addEventListener("beforeunload", saveCache);
    return () => window.removeEventListener("beforeunload", saveCache);
  }, []);

  return () => cache.current;
}

const loadCache = () => localStorage.getItem(CACHE_KEY_NAME);

const storeCache = (c: any) => localStorage.setItem(CACHE_KEY_NAME, c);
