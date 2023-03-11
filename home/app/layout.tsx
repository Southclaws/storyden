"use client";

import { CacheProvider } from "@chakra-ui/next-js";
import { ChakraProvider, extendTheme } from "@chakra-ui/react";
import { NextSeo } from "next-seo";
import localFont from "next/font/local";
import "./fonts.css";

const monasans = localFont({
  src: "./mona-sans.woff2",
  display: "swap",
});

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width" />
        <meta name="msapplication-TileColor" content="#303030" />
        <meta name="theme-color" content="#303030" />

        {/* Icons */}
        <link
          rel="apple-touch-icon"
          sizes="180x180"
          href="/apple-touch-icon.png"
        />
        <link
          rel="icon"
          type="image/png"
          sizes="32x32"
          href="/favicon-32x32.png"
        />
        <link
          rel="icon"
          type="image/png"
          sizes="16x16"
          href="/favicon-16x16.png"
        />
        <link rel="manifest" href="/site.webmanifest" />
        <link rel="mask-icon" href="/safari-pinned-tab.svg" color="#303030" />

        {/* SEO */}
        <NextSeo
          useAppDir={true}
          themeColor="#303030"
          title="Storyden: A forum for the modern age."
          // description="With a fresh new take on traditional bulletin board web forum software, Storyden is a modern, secure and extensible platform for building communities."
          description="Storyden is a platform for building communities. A modern take on oldschool bulletin board forums. Designed to be the community platform for the next era of internet culture."
          openGraph={{
            type: "website",
            locale: "en_GB",
            url: "https://www.storyden.org/",
            images: [
              {
                url: "https://www.storyden.org/social.png",
                width: 1728,
                height: 864,
                alt: "Storyden: A forum for the modern age.",
                type: "image/png",
              },
            ],
          }}
          twitter={{
            handle: "@Southclaws",
            site: "@Southclaws",
            cardType: "summary_large_image",
          }}
        />
      </head>
      <body>
        <CacheProvider>
          <ChakraProvider
            theme={extendTheme({
              fonts: {
                heading: "p22-mackinac-pro",
                body: monasans.style.fontFamily,
              },
            })}
          >
            {children}
          </ChakraProvider>
        </CacheProvider>
      </body>
    </html>
  );
}
