import type { Metadata, Viewport } from "next";
import { PropsWithChildren } from "react";

import { getColourAsHex } from "src/utils/colour";

import { inter, interDisplay } from "@/app/fonts";
import { serverEnvironment } from "@/config";
import { getSettings } from "@/lib/settings/settings-server";
import { filterAllowedThemeAssets } from "@/lib/theme/manifest";
import { getServerThemeManifest } from "@/lib/theme/theme-server";
import { getIconURL } from "@/utils/icon";

import "./global.css";

import { Providers } from "./providers";

const { API_ADDRESS, WEB_ADDRESS } = serverEnvironment();

export default async function RootLayout({ children }: PropsWithChildren) {
  const manifest = await getServerThemeManifest();
  const themeAssets = filterAllowedThemeAssets(manifest, {
    webAddress: WEB_ADDRESS,
    apiAddress: API_ADDRESS,
  });

  return (
    <html
      lang="en"
      className={`${inter.variable} ${interDisplay.variable} sd-app sd-app--root`}
      data-sd-theme-api="v1"
    >
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

        {themeAssets.css.map((href, i) => (
          // eslint-disable-next-line @next/next/no-css-tags
          <link
            key={`theme-css-${i}-${href}`}
            rel="stylesheet"
            href={href}
            data-sd-theme-asset="css"
          />
        ))}

        {themeAssets.scripts.map((src, i) => (
          <script
            key={`theme-script-${i}-${src}`}
            src={src}
            defer
            data-sd-theme-asset="script"
          />
        ))}
      </head>

      <body className="sd-app__body">
        <div className="sd-app__viewport">
          <Providers>{children}</Providers>
        </div>
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
