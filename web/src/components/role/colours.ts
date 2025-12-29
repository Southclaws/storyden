import chroma from "chroma-js";

export function badgeColourCSS(c: string) {
  const { bg, bo, fg } = badgeColours(c);

  return {
    "--colors-color-palette-fg": fg,
    "--colors-color-palette-border": bo,
    "--colors-color-palette-bg": bg,
    "--colors-color-palette-text": fg,
  } as React.CSSProperties;
}

export function badgeColours(c: string) {
  let colour;
  try {
    colour = chroma(c);
  } catch {
    // Default to green if colour is invalid or empty
    colour = chroma("#22c55e");
  }

  const bg = colour.brighten(2).desaturate(3).css();
  const bo = colour.darken(1).desaturate(1).alpha(0.2).css();
  const fg = colour.darken(1).saturate(2).css();

  return { bg, bo, fg };
}
