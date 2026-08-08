import chroma from "chroma-js";

import { getReadableTextColour } from "@/utils/colour";

export function categoryColourCSS(c: string) {
  const { bg, bo, fg } = categoryColours(c);

  return {
    "--colors-color-palette-fg": fg,
    "--colors-color-palette-border": bo,
    "--colors-color-palette-bg": bg,
  } as React.CSSProperties;
}

export function categoryColours(c: string) {
  const colour = chroma(c);

  const bg = colour.brighten(1).desaturate(1).css();
  const bo = colour.darken(0).desaturate(1).alpha(0.8).css();
  const fg = getReadableTextColour(bg);

  return { bg, bo, fg };
}
