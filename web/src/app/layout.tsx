import { PropsWithChildren } from "react";

import "./global.css";

import { Providers } from "./providers";

export default function RootLayout({ children }: PropsWithChildren) {
  return (
    <html lang="en">
      <head>
        {/*
            NOTE: This stylesheet is fully server-side rendered but it's not
            static because it uses data from the API to be generated. But we
            don't want this to require client-side render or CSS-in-JS.
        */}
        {/* eslint-disable-next-line @next/next/no-css-tags */}
        <link rel="stylesheet" href="/theme.css" />
      </head>

      <body suppressHydrationWarning={true}>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}

export const metadata = {
  manifest: "/manifest.json",
};
