import { createMDX } from 'fumadocs-mdx/next';

const withMDX = createMDX();

/** @type {import('next').NextConfig} */
const config = {
  reactStrictMode: true,
  async redirects() {
    return [
      {
        source: '/docs/operation/:path*',
        destination: '/docs/reference/:path*',
        permanent: true,
      },
    ];
  },
};

export default withMDX(config);
