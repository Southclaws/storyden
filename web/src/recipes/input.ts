import { defineRecipe } from "@pandacss/dev";

export const input = defineRecipe({
  className: "input",
  jsx: ["Input", "Field.Input"],
  base: {
    appearance: "none",
    background: "none",
    borderColor: "border.default",
    borderRadius: "l2",
    borderWidth: "1px",
    colorPalette: "accent",
    color: "fg.default",
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
      borderColor: "fg.error",
      _focus: {
        borderColor: "fg.error",
        boxShadow: "0 0 0 1px var(--colors-border-error)",
      },
    },
    _placeholder: {
      color: "fg.subtle",
      opacity: "full",
    },
  },
  defaultVariants: {
    size: "md",
    variant: "outline",
  },
  variants: {
    size: {
      "2xs": { px: "1.5", h: "6", minW: "7", fontSize: "xs" },
      xs: { px: "2", h: "6", minW: "8", fontSize: "xs", borderRadius: "sm" },
      sm: { px: "2.5", h: "8", minW: "9", fontSize: "sm" },
      md: { px: "3", h: "10", minW: "10", fontSize: "md" },
      lg: { px: "3.5", h: "11", minW: "11", fontSize: "md" },
      xl: { px: "4", h: "12", minW: "12", fontSize: "lg" },
      "2xl": { px: "4.5", h: "16", minW: "16", textStyle: "3xl" },
    },
    variant: {
      outline: {},
      ghost: {
        borderColor: "transparent",
        px: "0",
        _focus: {
          borderColor: "transparent",
          boxShadow: "0 0 0 0px var(--colors-color-palette-default)",
        },
      },
    },
  },
});
