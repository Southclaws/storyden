import { defineSemanticTokens } from "@pandacss/dev";

import blue from "@/theme/base/colours/blue";

import { accent } from "./accent";
import { bg } from "./bg";
import { border } from "./border";
import { conicGradient } from "./conic-gradient";
import { fg } from "./fg";
import { visibility } from "./visibility";

export const colours = defineSemanticTokens.colors({
  bg,
  fg,
  border,
  visibility,
  conicGradient,
  accent,

  blue: blue.semanticTokens,

  cardBackgroundGradient: {
    value: "linear-gradient(90deg, var(--colors-bg), transparent)",
  },
  backgroundGradientH: {
    value: "linear-gradient(90deg, var(--colors-bg), transparent)",
  },
  backgroundGradientV: {
    value: "linear-gradient(0deg, var(--colors-bg), transparent)",
  },
  "overflow-fade": {
    value: {
      _osLight:
        "linear-gradient(to right, rgb(from {colors.bg.subtle} r g b / 0) 0%, rgb(from {colors.bg.subtle} r g b / 1) 80%)",
      _osDark:
        "linear-gradient(to right, rgb(from {colors.bg.subtle} r g b / 0) 0%, rgb(from {colors.bg.subtle} r g b / 1) 80%)",
    },
  },
});
