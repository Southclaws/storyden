import chroma from "chroma-js";

import { badgeColourPalette } from "@/components/ui/badge/Badge.internal";
import { getReadableTextColour } from "@/utils/colour";

export function categoryColourCSS(c: string) {
  return badgeColourPalette(categoryColours(c));
}

export function categoryColours(c: string) {
  const colour = chroma(c);

  const background = colour.brighten(1).desaturate(1).css();
  const border = colour.darken(0).desaturate(1).alpha(0.8).css();
  const foreground = getReadableTextColour(background);

  return { background, border, foreground };
}
