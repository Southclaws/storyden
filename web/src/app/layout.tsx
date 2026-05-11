import type { Metadata, Viewport } from "next";
import { cookies } from "next/headers";
import { PropsWithChildren } from "react";

import { getColourAsHex } from "src/utils/colour";

import { inter, interDisplay } from "@/app/fonts";
import { serverEnvironment } from "@/config";
import { I18N_COOKIE_NAME, normalizeLocale } from "@/i18n/config";
import { getSettings } from "@/lib/settings/settings-server";
import { getIconURL } from "@/utils/icon";

import "./global.css";

import { Providers } from "./providers";

const { API_ADDRESS, WEB_ADDRESS } = serverEnvironment();

export default async function RootLayout({ children }: PropsWithChildren) {
  const cookieStore = await cookies();
  const locale = normalizeLocale(cookieStore.get(I18N_COOKIE_NAME)?.value);
  const browserConfig = JSON.stringify({
    API_ADDRESS,
    WEB_ADDRESS,
    source: "meta",
  });

  return (
    <html
      lang={locale}
      className={`${inter.variable} ${interDisplay.variable}`}
      suppressHydrationWarning
    >
      <head>
        {/*
          NOTE: Because the browser side does not support dynamic environment
          variables (obviously, it's a browser script) we hack around Next.js'
          build-time variables by embedding the public config in the document.
          This allows us to set the API/frontend addresses without rebuilding
          the entire app.
        */}
        <meta
          id="storyden-browser-config"
          name="storyden-browser-config"
          content={encodeURIComponent(browserConfig)}
        />

        {/*
            NOTE: This stylesheet is fully server-side rendered but it's not
            static because it uses data from the API to be generated. But we
            don't want this to require client-side render or CSS-in-JS.
        */}
        {/* eslint-disable-next-line @next/next/no-css-tags */}
        <link rel="stylesheet" href="/theme.css" />
      </head>

      <body>
        <Providers initialLocale={locale}>{children}</Providers>
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
