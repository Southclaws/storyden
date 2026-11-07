const isStandalone = process.env.NEXT_BUILD_STANDALONE === "true";

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: isStandalone ? "standalone" : undefined,

  reactStrictMode: true,
  reactCompiler: true,
  images: {
    loader: "custom",
    loaderFile: "./src/lib/asset/loader.js",
    unoptimized: true,
  },
  logging: {
    fetches: {
      fullUrl: true,
    },
  },
  experimental: {
    turbopackFileSystemCacheForDev: true,
  },

  async headers() {
    return [
      {
        source: "/sw.js",
        headers: [
          {
            key: "Content-Type",
            value: "application/javascript; charset=utf-8",
          },
          {
            key: "Cache-Control",
            value: "no-cache, no-store, must-revalidate",
          },
          {
            key: "Content-Security-Policy",
            value:
              "default-src 'none'; script-src 'self'; connect-src 'self' https: http:; img-src 'self' https: http:",
          },
        ],
      },
    ];
  },
};

module.exports = nextConfig;
