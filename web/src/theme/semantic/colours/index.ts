import { defineSemanticTokens } from "@pandacss/dev";

import amber from "@/theme/base/colours/amber";
import blue from "@/theme/base/colours/blue";
import green from "@/theme/base/colours/green";
import neutral from "@/theme/base/colours/neutral";
import orange from "@/theme/base/colours/orange";
import pink from "@/theme/base/colours/pink";
import red from "@/theme/base/colours/red";
import slate from "@/theme/base/colours/slate";
import tomato from "@/theme/base/colours/tomato";

import { accent } from "./accent";
import { background } from "./background";
import { border } from "./border";
import { conicGradient } from "./conic-gradient";
import { interactive } from "./interactive";
import { scrim } from "./scrim";
import { selection } from "./selection";
import { status } from "./status";
import { text } from "./text";
import { visibility } from "./visibility";

export const colours = defineSemanticTokens.colors({
  background,
  text,
  border,
  scrim,
  selection,
  interactive,
  status,
  visibility,
  conicGradient,
  accent,

  blue: blue.semanticTokens,
  green: green.semanticTokens,
  red: red.semanticTokens,
  amber: amber.semanticTokens,
  gray: neutral.semanticTokens,
  neutral: neutral.semanticTokens,
  orange: orange.semanticTokens,
  pink: pink.semanticTokens,
  slate: slate.semanticTokens,
  tomato: tomato.semanticTokens,

  cardBackgroundGradient: {
    value: "linear-gradient(90deg, {colors.background.surface}, transparent)",
  },
  backgroundGradientH: {
    value: "linear-gradient(90deg, {colors.background.surface}, transparent)",
  },
  backgroundGradientV: {
    value: "linear-gradient(0deg, {colors.background.surface}, transparent)",
  },
  "overflow-fade": {
    value: {
      _osLight:
        "linear-gradient(to right, rgb(from {colors.background.control} r g b / 0) 0%, rgb(from {colors.background.control} r g b / 1) 80%)",
      _osDark:
        "linear-gradient(to right, rgb(from {colors.background.control} r g b / 0) 0%, rgb(from {colors.background.control} r g b / 1) 80%)",
    },
  },
  "scroll-fade-top": {
    value: {
      _osLight:
        "linear-gradient(to bottom, {colors.background.canvas} 0%, transparent 100%)",
      _osDark:
        "linear-gradient(to bottom, {colors.background.canvas} 0%, transparent 100%)",
    },
  },
});
