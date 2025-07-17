const isStandalone = process.env.NEXT_BUILD_STANDALONE === "true";

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: isStandalone ? "standalone" : undefined,
  reactStrictMode: true,
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
    devtoolNewPanelUI: true,
    devtoolSegmentExplorer: true,
  },
};

module.exports = nextConfig;
