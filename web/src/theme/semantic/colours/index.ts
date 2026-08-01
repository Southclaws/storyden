import { defineSemanticTokens } from "@pandacss/dev";

import amber from "@/theme/base/colours/amber";
import blue from "@/theme/base/colours/blue";
import green from "@/theme/base/colours/green";
import red from "@/theme/base/colours/red";

import { accent } from "./accent";
import { background } from "./background";
import { badge } from "./badge";
import { border } from "./border";
import { conicGradient } from "./conic-gradient";
import { control } from "./control";
import { interactive } from "./interactive";
import { scrim } from "./scrim";
import { selection } from "./selection";
import { status } from "./status";
import { text } from "./text";
import { visibility } from "./visibility";

export const colours = defineSemanticTokens.colors({
  background,
  badge,
  text,
  border,
  control,
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
        "linear-gradient(to right, rgb(from {colors.control.background} r g b / 0) 0%, rgb(from {colors.control.background} r g b / 1) 80%)",
      _osDark:
        "linear-gradient(to right, rgb(from {colors.control.background} r g b / 0) 0%, rgb(from {colors.control.background} r g b / 1) 80%)",
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
