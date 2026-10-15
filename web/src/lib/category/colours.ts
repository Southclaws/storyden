import { badgeColourPalette } from "@/components/ui/badge/Badge.internal";
import { getReadableTextColour } from "@/utils/colour";
import { parseColour } from "@/utils/colour-core";

export function categoryColourCSS(c: string) {
  return badgeColourPalette(categoryColours(c));
}

export function categoryColours(c: string) {
  const colour = parseColour(c);

  const background = colour.brighten(1).desaturate(1).css();
  const border = colour.darken(0).desaturate(1).alpha(0.8).css();
  const foreground = getReadableTextColour(background);

  return { background, border, foreground };
}
