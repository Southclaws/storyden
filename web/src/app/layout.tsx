import type { Metadata } from "next";
import { PropsWithChildren } from "react";

import { getColourAsHex } from "src/utils/colour";
import { getInfo } from "src/utils/info";

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

export async function generateMetadata(): Promise<Metadata> {
  const info = await getInfo();

  const themeColour = getColourAsHex(info.accent_colour);
  const iconURL = `/api/v1/info/icon/512x512`;

  return {
    title: info.title,
    description: info.description,
    themeColor: themeColour,
    icons: {
      icon: iconURL,
      shortcut: iconURL,
      apple: iconURL,
    },
    appleWebApp: {
      capable: true,
      title: info.title,
      statusBarStyle: "black-translucent",
      startupImage: iconURL,
    },
  };
}
