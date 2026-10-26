const isStandalone = process.env.NEXT_BUILD_STANDALONE === "true";

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: isStandalone ? "standalone" : undefined,
  reactStrictMode: true,
  reactCompiler: true,
  images: {
    loader: "custom",
    loaderFile: "./src/lib/asset/loader.ts",
    deviceSizes: [320, 960, 1920],
    imageSizes: [],
  },
  logging: {
    fetches: {
      fullUrl: true,
    },
  },
  experimental: {
    turbopackFileSystemCacheForDev: true,
  },
};

module.exports = nextConfig;
