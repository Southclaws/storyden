import { defineTokens } from "@pandacss/dev";

import { animations } from "./animations";
import { blurs } from "./blurs";
import { borders } from "./borders";
import { colours } from "./colours";
import { durations } from "./durations";
import { easings } from "./easings";
import { radii } from "./radii";
import { shadows } from "./shadows";
import { sizes } from "./sizes";
import { spacing } from "./spacing";
import {
  fontSizes,
  fontWeights,
  fonts,
  letterSpacings,
  lineHeights,
} from "./typography";
import { zIndex } from "./z-index";

export const tokens = defineTokens({
  animations,
  blurs,
  borders,
  colors: colours,
  durations,
  easings,
  fonts,
  fontSizes,
  fontWeights,
  letterSpacings,
  lineHeights,
  radii,
  sizes,
  spacing,
  zIndex,
  shadows,
});
