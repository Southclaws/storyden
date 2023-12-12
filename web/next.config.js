/** @type {import('next').NextConfig} */

const isStandalone = process.env.NEXT_BUILD_STANDALONE === "true";

const nextConfig = {
  output: isStandalone ? "standalone" : undefined,
  reactStrictMode: true,
  swcMinify: true,
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${
          process.env.NEXT_PUBLIC_API_ADDRESS ?? "http://localhost:8000"
        }/api/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;
