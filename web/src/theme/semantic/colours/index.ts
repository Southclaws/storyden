import { defineSemanticTokens } from "@pandacss/dev";

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

  cardBackgroundGradient: {
    value: "linear-gradient(90deg, var(--colors-bg), transparent)",
  },
  backgroundGradientH: {
    value: "linear-gradient(90deg, var(--colors-bg), transparent)",
  },
  backgroundGradientV: {
    value: "linear-gradient(0deg, var(--colors-bg), transparent)",
  },
});
