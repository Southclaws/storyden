import { defineSemanticTokens } from "@pandacss/dev";

import { colours } from "./colours";

export const semanticTokens = defineSemanticTokens({
  colors: colours,
  radii: {
    l1: { value: "{radii.xs}" },
    l2: { value: "{radii.sm}" },
    l3: { value: "{radii.md}" },
  },
  fonts: {
    body: { value: "{fonts.inter}" },
    heading: { value: "{fonts.interDisplay}" },
  },
  blurs: {
    frosted: { value: "10px" },
  },
  opacity: {
    0: { value: "0" },
    1: { value: "0.1" },
    2: { value: "0.2" },
    3: { value: "0.3" },
    4: { value: "0.4" },
    5: { value: "0.5" },
    6: { value: "0.6" },
    7: { value: "0.7" },
    8: { value: "0.8" },
    9: { value: "0.9" },
    full: { value: "1" },
  },
  borderWidths: {
    none: { value: "0" },
    hairline: { value: "0.5px" },
    thin: { value: "1px" },
    medium: { value: "2px" },
    thick: { value: "3px" },
  },
  sizes: {
    prose: { value: "65ch" },
  },
  spacing: {
    safeBottom: { value: "env(safe-area-inset-bottom)" },
    scrollGutter: { value: "var(--spacing-2)" },
  },
});
