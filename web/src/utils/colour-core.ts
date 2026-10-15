import { css } from "chroma-js/src/io/css/index.js";
import { hex } from "chroma-js/src/io/hex/index.js";
import "chroma-js/src/io/hsl/index.js";
import "chroma-js/src/io/named/index.js";
import "chroma-js/src/io/rgb/index.js";
import "chroma-js/src/ops/alpha.js";
import "chroma-js/src/ops/darken.js";
import "chroma-js/src/ops/saturate.js";

// Chroma's main entry registers every generator, interpolator and colour space.
// Keep the existing conversions and operators through its exported subpaths.
export { oklch } from "chroma-js/src/io/oklch/index.js";
export { default as contrast } from "chroma-js/src/utils/contrast.js";

export function parseColour(value: string) {
  return /^#?[\da-f]{3,8}$/i.test(value) ? hex(value) : css(value);
}
