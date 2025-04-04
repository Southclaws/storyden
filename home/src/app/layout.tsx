import { RootProvider } from "fumadocs-ui/provider";
import "fumadocs-ui/style.css";
import type { ReactNode } from "react";

import { joie, worksans, hedvig, intelone } from "@/fonts";
import "./globals.css";

import Script from "next/script";
import { cx } from "@/styled-system/css";

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <html
      lang="en"
      className={cx(
        joie.variable,
        worksans.variable,
        hedvig.variable,
        intelone.variable
      )}
      suppressHydrationWarning
    >
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
            defaultTheme: "light",
            enableSystem: false,
          }}
        >
          {children}
        </RootProvider>
      </body>
    </html>
  );
}
