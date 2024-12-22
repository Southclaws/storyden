import type { Metadata, Viewport } from "next";
import { PropsWithChildren } from "react";

import { getColourAsHex } from "src/utils/colour";

import { inter, interDisplay } from "@/app/fonts";
import { getSettings } from "@/lib/settings/settings-server";
import { getIconURL } from "@/utils/icon";

import "./global.css";

import { Providers } from "./providers";

export const dynamic = "force-dynamic";

const API_ADDRESS =
  global.process.env["NEXT_PUBLIC_API_ADDRESS"] ?? "http://localhost:8000";
const WEB_ADDRESS =
  global.process.env["NEXT_PUBLIC_WEB_ADDRESS"] ?? "http://localhost:3000";

export default async function RootLayout({ children }: PropsWithChildren) {
  return (
    <html lang="en" className={`${inter.variable} ${interDisplay.variable}`}>
      <head>
        {/*
          NOTE: Because the browser side does not support dynamic environment
          variables (obviously, it's a browser script) we hack around Next.js'
          build-time variables by providing a direct reference to these inside
          the window object. This allows us to set the API/frontend addresses
          without rebuilding the entire app.
        */}
        <script>{`
          window.__storyden__ = {"API_ADDRESS":"${API_ADDRESS}", "WEB_ADDRESS":"${WEB_ADDRESS}", "source": "script"};
          console.log("set up window config", window.__storyden__);
        `}</script>

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
  const settings = await getSettings();

  const themeColour = getColourAsHex(settings.accent_colour);

  return {
    themeColor: themeColour,
    colorScheme: "only light",
  };
}

export async function generateMetadata(): Promise<Metadata> {
  const settings = await getSettings();

  const iconURL = getIconURL("512x512");

  const canonical = WEB_ADDRESS;

  // TODO: Add another settings field for this.
  const title = `${settings.title} | ${settings.description}`;

  return {
    manifest: "/manifest.json",
    metadataBase: new URL(canonical),
    title: title,
    description: settings.description,
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
