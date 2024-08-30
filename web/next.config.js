/** @type {import('next').NextConfig} */

const isStandalone = process.env.NEXT_BUILD_STANDALONE === "true";

const API_ADDRESS =
  process.env["NEXT_PUBLIC_API_ADDRESS"] ?? "http://localhost:8000";

const apiURL = new URL(API_ADDRESS);

const nextConfig = {
  output: isStandalone ? "standalone" : undefined,
  reactStrictMode: true,
  swcMinify: true,
  images: {
    remotePatterns: [
      {
        protocol: apiURL.protocol.replace(":", ""),
        hostname: apiURL.hostname,
        port: apiURL.port,
      },
    ],
  },
};

module.exports = nextConfig;
