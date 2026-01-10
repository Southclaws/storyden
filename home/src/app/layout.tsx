import { RootProvider } from "fumadocs-ui/provider/next";
import "fumadocs-ui/style.css";
import type { ReactNode } from "react";

import { joie, worksans, hedvig, intelone, gorton } from "@/fonts";
import "./globals.css";

import Script from "next/script";
import { cx } from "@/styled-system/css";
import { Metadata, Viewport } from "next";

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <html
      lang="en"
      className={cx(
        joie.variable,
        worksans.variable,
        hedvig.variable,
        intelone.variable,
        gorton.variable
      )}
      suppressHydrationWarning
    >
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width" />
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
          sizes="96x96"
          href="/favicon-96x96.png"
        />
        <link rel="manifest" href="/site.webmanifest" />
        <Script type="text/javascript">
          {`(function(c,l,a,r,i,t,y){
            c[a]=c[a]||function(){(c[a].q=c[a].q||[]).push(arguments)};
            t=l.createElement(r);t.async=1;t.src="https://www.clarity.ms/tag/"+i;
            y=l.getElementsByTagName(r)[0];y.parentNode.insertBefore(t,y);
        })(window, document, "clarity", "script", "obgioniw76");`}
        </Script>
      </head>
      <body
        style={{
          display: "flex",
          flexDirection: "column",
          minHeight: "100vh",
        }}
      >
        <RootProvider
          theme={{
            defaultTheme: "system",
            enableSystem: true,
          }}
        >
          {children}
        </RootProvider>
      </body>
    </html>
  );
}

export const metadata: Metadata = {
  metadataBase: new URL("https://www.storyden.org"),
  title: "Storyden: A forum for the modern age.",
  description:
    "Storyden is a platform for building communities. A modern take on oldschool bulletin board forums. Designed to be the community platform for the next era of internet culture.",
  openGraph: {
    type: "website",
    locale: "en_GB",
    url: "https://www.storyden.org/",
    images: [
      {
        url: "https://www.storyden.org/opengraph-1280-640.png",
        width: 1280,
        height: 640,
        alt: "Storyden: A forum for the modern age.",
        type: "image/png",
      },
    ],
  },

  twitter: {
    creator: "@Southclaws",
    site: "@Southclaws",
    card: "summary_large_image",
  },

  alternates: {
    types: {
      "application/rss+xml": [
        {
          title: "Storyden Blog",
          url: "https://www.storyden.org/rss.xml",
        },
      ],
    },
  },
};

export const viewport: Viewport = {
  themeColor: "#d8dbcd",
};
