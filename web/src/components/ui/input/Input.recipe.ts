import { defineRecipe } from "@pandacss/dev";

export const inputBase = {
  appearance: "none",
  background: "background.control",
  borderColor: "border.default",
  borderRadius: "sm",
  borderWidth: "1px",
  colorPalette: "accent",
  color: "text.default",
  outline: 0,
  position: "relative",
  transitionDuration: "normal",
  transitionProperty: "box-shadow, border-color",
  transitionTimingFunction: "default",
  width: "full",
  _disabled: {
    opacity: 0.4,
    cursor: "not-allowed",
  },
  _focus: {
    outline: "0",
    borderColor: "none",
    boxShadow: "none",
  },
  _invalid: {
    borderColor: "status.danger.content",
    _focus: {
      borderColor: "status.danger.content",
      boxShadow: "0 0 0 1px var(--colors-border-error)",
    },
  },
  _placeholder: {
    color: "text.muted",
    opacity: "full",
  },
} as const;

export const inputSizeVariants = {
  sm: { h: "6", minW: "8", fontSize: "xs", borderRadius: "sm" },
  md: { h: "8", minW: "9", fontSize: "sm" },
  lg: { h: "10", minW: "10", fontSize: "md" },
} as const;

export const inputVariantVariants = {
  outline: {},
  ghost: {
    borderColor: "transparent",
    px: "0",
    _focus: {
      borderColor: "transparent",
      boxShadow: "0 0 0 0px var(--colors-color-palette-default)",
    },
  },
  inset: {
    background: "background.controlInset",
    borderColor: "transparent",
    px: "0",
    _focus: {
      borderColor: "transparent",
      boxShadow: "0 0 0 0px var(--colors-color-palette-default)",
    },
  },
} as const;

export const input = defineRecipe({
  className: "input",
  jsx: ["Input", "Field.Input"],
  base: inputBase,
  defaultVariants: {
    size: "sm",
    variant: "outline",
  },
  variants: {
    size: inputSizeVariants,
    variant: inputVariantVariants,
  },
  compoundVariants: [
    { size: "sm", variant: "outline", css: { px: "2" } },
    { size: "md", variant: "outline", css: { px: "2.5" } },
    { size: "lg", variant: "outline", css: { px: "3" } },
  ],
});
