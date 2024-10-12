import chroma from "chroma-js";

export function categoryColourCSS(c: string) {
  const { bg, bo, fg } = categoryColours(c);

  return {
    "--colors-color-palette-text": fg,
    "--colors-color-palette-muted": bo,
    "--colors-color-palette": bg,
  } as React.CSSProperties;
}

export function categoryColours(c: string) {
  const colour = chroma(c);

  const bg = colour.brighten(3).desaturate(2).css();
  const bo = colour.darken(0).desaturate(1).alpha(0.8).css();
  const fg = colour.darken(1).saturate(2).css();

  return { bg, bo, fg };
}
