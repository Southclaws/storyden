import { defineRecipe } from "@pandacss/dev";

export const button = defineRecipe({
  className: "button",
  jsx: ["Button", "IconButton", "SubmitButton"],
  base: {
    alignItems: "center",
    appearance: "none",
    borderWidth: "1px",
    borderColor: "transparent",
    borderRadius: "l2",
    boxShadow: "xs",
    cursor: "pointer",
    display: "inline-flex",
    fontWeight: "medium",
    minWidth: "0",
    justifyContent: "center",
    outline: "none",
    position: "relative",
    transitionDuration: "normal",
    transitionProperty: "background, border-color, color, box-shadow",
    transitionTimingFunction: "default",
    userSelect: "none",
    verticalAlign: "middle",
    whiteSpace: "nowrap",
    _hidden: {
      display: "none",
    },
  },
  defaultVariants: {
    variant: "solid",
    size: "md",
  },
  variants: {
    variant: {
      solid: {
        background: "colorPalette.default",
        color: "bg.default",
        colorPalette: "accent",
        _hover: {
          background: "colorPalette.default/90",
          color: "bg.default",
        },
        _active: {
          background: "colorPalette.default/90",
        },
        _focusVisible: {
          boxShadow: "0 0 0 3px var(--colors-border-focus)",
        },
        _disabled: {
          color: "fg.disabled",
          background: "bg.disabled",
          cursor: "not-allowed",
          _hover: {
            color: "fg.disabled",
            background: "bg.disabled",
          },
        },
      },
      outline: {
        borderWidth: "1px",
        borderColor: "border.default",
        background: "bg.default",
        color: "fg.subtle",
        colorPalette: "gray",
        _hover: {
          borderColor: "border.strong",
          background: "bg.subtle",
          color: "fg.default",
        },
        _disabled: {
          borderColor: "border.disabled",
          color: "fg.disabled",
          cursor: "not-allowed",
          _hover: {
            background: "transparent",
            borderColor: "border.disabled",
            color: "fg.disabled",
          },
        },
        _focusVisible: {
          boxShadow: "0 0 0 3px var(--colors-border-focus)",
        },
        _selected: {
          background: "accent.default",
          borderColor: "accent.default",
          color: "accent.fg",
          _hover: {
            background: "accent.default/80",
            borderColor: "accent.default/80",
          },
        },
      },
      ghost: {
        color: "fg.subtle",
        colorPalette: "accent",
        _hover: {
          background: "bg.subtle",
          color: "fg.default",
        },
        _selected: {
          background: "colorPalette.muted",
        },
        _disabled: {
          color: "fg.disabled",
          cursor: "not-allowed",
          _hover: {
            background: "transparent",
            color: "fg.disabled",
          },
        },
        _focusVisible: {
          boxShadow: "0 0 0 3px var(--colors-border-focus)",
        },
      },
      link: {
        verticalAlign: "baseline",
        _disabled: {
          color: "border.disabled",
          cursor: "not-allowed",
          _hover: {
            color: "border.disabled",
          },
        },
        height: "auto!",
        px: "0!",
        minW: "0!",
      },
      subtle: {
        colorPalette: "accent",
        background: "bg.muted",
        borderColor: "border.subtle",
        color: "fg.subtle",
        _hover: {
          borderColor: "border.default",
          background: "bg.subtle",
          color: "fg.default",
        },
        _focusVisible: {
          boxShadow: "0 0 0 3px var(--colors-border-focus)",
        },
        _disabled: {
          background: "bg.disabled",
          color: "fg.disabled",
          cursor: "not-allowed",
          _hover: {
            background: "bg.disabled",
            color: "fg.disabled",
          },
        },
      },
    },
    size: {
      xs: {
        h: "6",
        minW: "6",
        textStyle: "xs",
        borderRadius: "sm",
        px: "2",
        gap: "2",
        "& svg": {
          flexShrink: "0",
          fontSize: "sm",
          width: "4",
          height: "4",
        },
      },
      sm: {
        h: "8",
        minW: "8",
        textStyle: "sm",
        px: "3",
        gap: "2",
        "& svg": {
          width: "4",
          height: "4",
        },
      },
      md: {
        h: "9",
        minW: "9",
        textStyle: "sm",
        px: "3.5",
        gap: "2",
        "& svg": {
          width: "5",
          height: "5",
        },
      },
      lg: {
        h: "10",
        minW: "10",
        textStyle: "md",
        px: "4",
        gap: "2",
        "& svg": {
          width: "5",
          height: "5",
        },
      },
      xl: {
        h: "12",
        minW: "12",
        textStyle: "md",
        px: "5",
        gap: "2.5",
        "& svg": {
          width: "5",
          height: "5",
        },
      },
      "2xl": {
        h: "16",
        minW: "16",
        textStyle: "lg",
        px: "7",
        gap: "3",
        "& svg": {
          width: "6",
          height: "6",
        },
      },
    },
  },
});
