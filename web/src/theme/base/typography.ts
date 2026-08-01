import { defineTokens } from "@pandacss/dev";

export const fontSizes = defineTokens.fontSizes({
  "2xs": { value: "0.5rem" },
  xs: { value: "0.75rem" },
  sm: { value: "0.875rem" },
  md: { value: "1rem" },
  lg: { value: "1.125rem" },
  xl: { value: "1.25rem" },
  "2xl": { value: "1.5rem" },
});

export const fontWeights = defineTokens.fontWeights({
  thin: { value: "100" },
  extralight: { value: "200" },
  light: { value: "300" },
  normal: { value: "400" },
  medium: { value: "500" },
  semibold: { value: "600" },
  bold: { value: "700" },
  extrabold: { value: "800" },
  black: { value: "900" },
});

export const letterSpacings = defineTokens.letterSpacings({
  tighter: { value: "-0.05em" },
  tight: { value: "-0.025em" },
  normal: { value: "0em" },
  wide: { value: "0.025em" },
  wider: { value: "0.05em" },
  widest: { value: "0.1em" },
});

export const lineHeights = defineTokens.lineHeights({
  none: { value: "1" },
  tight: { value: "1.25" },
  normal: { value: "1.5" },
  relaxed: { value: "1.75" },
  loose: { value: "2" },
});

export const fonts = defineTokens.fonts({
  inter: { value: "var(--font-inter)" },
  interDisplay: { value: "var(--font-inter-display)" },
  mono: {
    value: [
      "ui-monospace",
      "SFMono-Regular",
      "Menlo",
      "Monaco",
      "Consolas",
      '"Liberation Mono"',
      '"Courier New"',
      "monospace",
    ],
  },
});
