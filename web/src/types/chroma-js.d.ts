// @types/chroma-js only declares the main entry, not the package's exported subpaths.
declare module "chroma-js/src/io/css/index.js" {
  import type { Color } from "chroma-js";
  export function css(value: string): Color;
}

declare module "chroma-js/src/io/hex/index.js" {
  import type { Color } from "chroma-js";
  export function hex(value: string): Color;
}

declare module "chroma-js/src/io/oklch/index.js" {
  import type { Color } from "chroma-js";
  export function oklch(lightness: number, chroma: number, hue: number): Color;
}

declare module "chroma-js/src/utils/contrast.js" {
  import type { Color } from "chroma-js";
  export default function contrast(
    a: string | Color,
    b: string | Color,
  ): number;
}
