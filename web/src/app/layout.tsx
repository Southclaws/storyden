import type { Metadata, Viewport } from "next";
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

      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}

export async function generateViewport(): Promise<Viewport> {
  const info = await getInfo();

  const themeColour = getColourAsHex(info.accent_colour);

  return {
    themeColor: themeColour,
    colorScheme: "only light",
  };
}

export async function generateMetadata(): Promise<Metadata> {
  const info = await getInfo();

  const iconURL = `/api/v1/info/icon/512x512`;

  const canonical =
    process.env["NEXT_PUBLIC_WEB_ADDRESS"] ?? "http://localhost:3000";

  // TODO: Add another settings field for this.
  const title = `${info.title} | ${info.description}`;

  return {
    manifest: "/manifest.json",
    metadataBase: new URL(canonical),
    title: title,
    description: info.description,
    icons: {
      icon: iconURL,
      shortcut: iconURL,
      apple: iconURL,
    },
    appleWebApp: {
      capable: true,
      title: title,
      statusBarStyle: "black-translucent",
      startupImage: iconURL,
    },
  };
}
