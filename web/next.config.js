const path = require("path");
const { InjectManifest } = require("workbox-webpack-plugin");

const isStandalone = process.env.NEXT_BUILD_STANDALONE === "true";

const API_ADDRESS =
  process.env["NEXT_PUBLIC_API_ADDRESS"] ?? "http://localhost:8000";

const apiURL = new URL(API_ADDRESS);

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: isStandalone ? "standalone" : undefined,
  reactStrictMode: true,
  images: {
    remotePatterns: [
      {
        protocol: apiURL.protocol.replace(":", ""),
        hostname: apiURL.hostname,
        port: apiURL.port,
      },
    ],
  },
  webpack: (config, options) => {
    if (options.isServer) {
      return config;
    }

    const dev = options.dev;

    const swDest = path.join(options.dir, "public", "sw.js");

    const workboxPlugin = new InjectManifest({
      compileSrc: true,

      dontCacheBustURLsMatching: /^\/_next\/static\/.*/iu,
      maximumFileSizeToCacheInBytes: 1024 * 1024 * 20,
      swSrc: "./src/worker/worker.ts",
      swDest,
      modifyURLPrefix: {
        "/": "/_next/static/",
      },

      // In dev, exclude everything.
      // This avoids irrelevant warnings about chunks being too large for caching.
      // In non-dev, use the default `exclude` option, don't override.
      // ...(dev ? { exclude: [/./] } : undefined),
    });

    if (options.dev) {
      // Suppress the "InjectManifest has been called multiple times" warning by reaching into
      // the private properties of the plugin and making sure it never ends up in the state
      // where it makes that warning.
      // https://github.com/GoogleChrome/workbox/blob/v6/packages/workbox-webpack-plugin/src/inject-manifest.ts#L260-L282
      Object.defineProperty(workboxPlugin, "alreadyCalled", {
        get() {
          return false;
        },
        set() {
          // do nothing; the internals try to set it to true, which then results in a warning
          // on the next run of webpack.
        },
      });
    }

    const newConfig = {
      ...config,
      plugins: [...config.plugins, workboxPlugin],
    };

    return newConfig;
  },
};

module.exports = nextConfig;
