import type { Metadata } from "next";
import { PropsWithChildren } from "react";

import { getIconGetKey } from "src/api/openapi/misc";
import { Default } from "src/layouts/Default";
import { getColourAsHex } from "src/utils/colour";
import { getInfo } from "src/utils/info";

export default function Layout({ children }: PropsWithChildren) {
  return <Default>{children}</Default>;
}

export async function generateMetadata(): Promise<Metadata> {
  const info = await getInfo();

  const themeColour = getColourAsHex(info.accent_colour);
  const iconLarge = getIconGetKey("512x512")[0];
  const iconURL = `/api${iconLarge}`;

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
