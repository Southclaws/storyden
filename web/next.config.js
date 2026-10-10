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
  async rewrites() {
    return {
      beforeFiles: [
        { source: "/t/:slug.md", destination: "/markdown/t/:slug" },
        { source: "/l/:slug(.*).md", destination: "/markdown/l/:slug" },
      ],
    };
  },
};

module.exports = nextConfig;
