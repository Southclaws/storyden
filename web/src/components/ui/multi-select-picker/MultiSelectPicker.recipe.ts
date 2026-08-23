import { defineRecipe } from "@pandacss/dev";

export const multiSelectPicker = defineRecipe({
  className: "multi-select-picker",
  base: {
    alignItems: "center",
    appearance: "none",
    background: "background.control",
    borderColor: "border.default",
    borderRadius: "control",
    borderWidth: "thin",
    color: "text.default",
    cursor: "text",
    display: "flex",
    gap: "1",
    minWidth: "0",
    outline: "none",
    position: "relative",
    transitionDuration: "normal",
    transitionProperty: "border-color, box-shadow, background",
    transitionTimingFunction: "default",
    width: "full",
    _hover: {
      borderColor: "border.strong",
    },
    _focusWithin: {
      borderColor: "colorPalette.text",
      boxShadow: "0 0 0 1px var(--colors-color-palette-text)",
    },
    _disabled: {
      background: "background.controlDisabled",
      borderColor: "border.disabled",
      color: "text.disabled",
      cursor: "not-allowed",
      _hover: {
        borderColor: "border.disabled",
      },
    },
  },
  defaultVariants: {
    size: "sm",
  },
  variants: {
    size: {
      sm: {
        minHeight: "6",
        paddingBlock: "0",
        paddingLeft: "0.5",
        paddingRight: "6",
      },
      md: {
        minHeight: "8",
        paddingBlock: "0",
        paddingLeft: "1",
        paddingRight: "8",
      },
      lg: {
        minHeight: "10",
        paddingBlock: "0",
        paddingLeft: "1.5",
        paddingRight: "10",
      },
    },
  },
});
