/// <reference lib="es2017" />
/// <reference lib="WebWorker" />
import {
  RouteHandlerCallback,
  RouteHandlerCallbackOptions,
  clientsClaim,
  skipWaiting,
} from "workbox-core";
import { ExpirationPlugin } from "workbox-expiration";
import { cleanupOutdatedCaches, precacheAndRoute } from "workbox-precaching";
import { PrecacheFallbackPlugin } from "workbox-precaching";
import { offlineFallback, staticResourceCache } from "workbox-recipes";
import { registerRoute, setDefaultHandler } from "workbox-routing";
import {
  NetworkFirst,
  NetworkOnly,
  StaleWhileRevalidate,
} from "workbox-strategies";

import { API_ADDRESS } from "@/config";

declare const self: ServiceWorkerGlobalScope;

self.__WB_DISABLE_DEV_LOGS = false;

console.log("Loading: worker.ts");

setDefaultHandler(new NetworkOnly());

offlineFallback();

skipWaiting();
clientsClaim();
// cleanupOutdatedCaches();
precacheAndRoute(self.__WB_MANIFEST);

staticResourceCache({
  cacheName: "static-resources",
  plugins: [
    new ExpirationPlugin({
      maxEntries: 120,
      maxAgeSeconds: 24 * 60 * 60, // 24 hours
      purgeOnQuotaError: true,
    }),
  ],
});

// imageCache()

// setDefaultHandler((options: RouteHandlerCallbackOptions): Promise<Response> => {
//   console.log("setDefaultHandler", options);

//   return fetch(options.request);
// });

registerRoute(
  /\/_next\/data\/.+\/.+\.json$/i,
  new StaleWhileRevalidate({
    matchOptions: {
      ignoreVary: true,
    },
    cacheName: "next-data",
    plugins: [
      new ExpirationPlugin({
        maxEntries: 120,
        maxAgeSeconds: 24 * 60 * 60, // 24 hours
        purgeOnQuotaError: true,
      }),
    ],
  }),
);

registerRoute(
  ({ request }) => request.mode === "navigate",
  new NetworkFirst({
    cacheName: "navigate",
    matchOptions: {
      ignoreVary: true,
    },
  }),
);

registerRoute(
  ({ url }) => {
    const origin = url.origin;
    const path = url.pathname;

    // Match only Storyden API requests.
    const originMatch = origin === API_ADDRESS;
    const pathMatch = path.startsWith("/api");

    const match = originMatch && pathMatch;

    return match;
  },

  new StaleWhileRevalidate({
    cacheName: "api",
    plugins: [
      new ExpirationPlugin({
        maxEntries: 120,
        maxAgeSeconds: 60 * 60,
        purgeOnQuotaError: true,
      }),
    ],
  }),
);

addEventListener("message", (event) => {
  if (event.data && event.data.type === "SKIP_WAITING") {
    self.skipWaiting();
  }
});
