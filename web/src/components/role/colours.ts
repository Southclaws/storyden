import chroma from "chroma-js";

import { getReadableTextColour } from "@/utils/colour";

import { badgeColourPalette } from "../ui/badge/Badge.internal";

export function badgeColourCSS(c: string) {
  return badgeColourPalette(badgeColours(c));
}

export function badgeColours(c: string) {
  let colour;
  try {
    colour = chroma(c);
  } catch {
    // Default to green if colour is invalid or empty
    colour = chroma("#22c55e");
  }

  const background = colour.brighten(2).desaturate(3).css();
  const border = colour.darken(1).desaturate(1).alpha(0.2).css();
  const foreground = getReadableTextColour(background);

  return { background, border, foreground };
}
