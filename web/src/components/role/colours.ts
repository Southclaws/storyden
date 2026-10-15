import { getReadableTextColour } from "@/utils/colour";
import { parseColour } from "@/utils/colour-core";

import { badgeColourPalette } from "../ui/badge/Badge.internal";

export function badgeColourCSS(c: string) {
  return badgeColourPalette(badgeColours(c));
}

export function badgeColours(c: string) {
  let colour;
  try {
    colour = parseColour(c);
  } catch {
    // Default to green if colour is invalid or empty
    colour = parseColour("#22c55e");
  }

  const background = colour.brighten(2).desaturate(3).css();
  const border = colour.darken(1).desaturate(1).alpha(0.2).css();
  const foreground = getReadableTextColour(background);

  return { background, border, foreground };
}
