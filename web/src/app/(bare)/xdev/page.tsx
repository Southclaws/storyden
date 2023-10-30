import { WEB_ADDRESS } from "src/config";
import { Palette } from "src/screens/xdev/xdev";
import { getColourVariants } from "src/utils/colour";
import { getInfo } from "src/utils/info";

export default async function Page() {
  const theme = await fetch(`${WEB_ADDRESS}/theme.css`);
  const info = await getInfo();
  const themeText = await theme.text();

  const colours = getColourVariants(info.accent_colour);

  return (
    <Palette
      accent_colour={info.accent_colour}
      colours={colours}
      info={info}
      theme={themeText}
    />
  );
}
