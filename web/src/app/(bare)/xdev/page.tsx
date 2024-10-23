import { WEB_ADDRESS } from "src/config";
import { Palette } from "src/screens/xdev/xdev";
import { getColourVariants } from "src/utils/colour";

import { getSettings } from "@/lib/settings/settings-server";

export const dynamic = "force-dynamic";

export default async function Page() {
  const theme = await fetch(`${WEB_ADDRESS}/theme.css`);
  const settings = await getSettings();
  const themeText = await theme.text();

  const colours = getColourVariants(settings.accent_colour);

  return (
    <Palette
      accent_colour={settings.accent_colour}
      colours={colours}
      info={settings}
      theme={themeText}
    />
  );
}
