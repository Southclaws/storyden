import { defineRecipe } from "@pandacss/dev";

export const button = defineRecipe({
  className: "button",
  jsx: ["Button", "IconButton", "SubmitButton"],
  base: {
    alignItems: "center",
    appearance: "none",
    borderRadius: "sm",
    cursor: "pointer",
    display: "inline-flex",
    fontWeight: "semibold",
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
    "& .button__loading-content": {
      alignItems: "center",
      display: "inline-flex",
      gap: "inherit",
      opacity: "0",
    },
    "& .button__loading-indicator": {
      alignItems: "center",
      display: "flex",
      inset: "0",
      justifyContent: "center",
      position: "absolute",
    },
    _hidden: {
      display: "none",
    },
  },
  defaultVariants: {
    variant: "subtle",
    size: "sm",
  },
  variants: {
    size: {
      sm: {
        h: "6",
        minW: "6",
        fontSize: "xs",
        lineHeight: "1.125rem",
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
      md: {
        h: "8",
        minW: "8",
        fontSize: "sm",
        lineHeight: "1.25rem",
        px: "2",
        gap: "2",
        "& svg": {
          width: "4",
          height: "4",
        },
      },
      lg: {
        h: "10",
        minW: "10",
        fontSize: "sm",
        lineHeight: "1.25rem",
        px: "4",
        gap: "2",
        "& svg": {
          width: "5",
          height: "5",
        },
      },
    },
    variant: {
      solid: {
        background: "colorPalette.default",
        color: "colorPalette.fg",
        colorPalette: "accent",
        _hover: {
          background: "colorPalette.default/80",
          color: "colorPalette.fg",
        },
        _focusVisible: {
          outline: "2px solid",
          outlineColor: "colorPalette.default",
          outlineOffset: "2px",
        },
        _disabled: {
          color: "text.disabled",
          background: "background.controlDisabled",
          cursor: "not-allowed",
          _hover: {
            color: "text.disabled",
            background: "background.controlDisabled",
          },
        },
      },
      outline: {
        borderWidth: "1px",
        borderColor: "border.muted",
        color: "text.muted",
        colorPalette: "gray",
        _hover: {
          background: "colorPalette.subtle",
          color: "colorPalette.default",
        },
        _disabled: {
          borderColor: "border.disabled",
          color: "text.disabled",
          cursor: "not-allowed",
          _hover: {
            background: "transparent",
            borderColor: "border.disabled",
            color: "text.disabled",
          },
        },
        _focusVisible: {
          outline: "2px solid",
          outlineColor: "colorPalette.default",
          outlineOffset: "2px",
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
        color: "text.muted",
        colorPalette: "accent",
        _hover: {
          background: "colorPalette.subtle",
          color: "text.default",
        },
        _selected: {
          background: "colorPalette.muted",
        },
        _disabled: {
          color: "text.disabled",
          cursor: "not-allowed",
          _hover: {
            background: "transparent",
            color: "text.disabled",
          },
        },
        _focusVisible: {
          outline: "2px solid",
          outlineColor: "colorPalette.default",
          outlineOffset: "2px",
        },
      },
      subtle: {
        colorPalette: "accent",
        background: "background.controlSubtle/80",
        color: "text.muted",
        backdropBlur: "sm",
        backdropFilter: "auto",
        _hover: {
          background: "colorPalette.subtle",
          color: "text.default",
        },
        _focusVisible: {
          outline: "2px solid",
          outlineColor: "colorPalette.default",
          outlineOffset: "2px",
        },
        _disabled: {
          background: "background.controlDisabled",
          color: "text.disabled",
          cursor: "not-allowed",
          _hover: {
            background: "background.controlDisabled",
            color: "text.disabled",
          },
        },
      },
      plain: {
        background: "transparent",
        color: "text.muted",
        gap: "1",
        h: "auto",
        minW: "0",
        p: "0",
        _hover: {
          background: "transparent",
          color: "text.default",
        },
        _focusVisible: {
          outline: "2px solid",
          outlineColor: "border.default",
          outlineOffset: "2px",
        },
        _disabled: {
          background: "transparent",
          color: "text.disabled",
          cursor: "not-allowed",
          _hover: {
            background: "transparent",
            color: "text.disabled",
          },
        },
        _after: {
          content: '\"\"',
          height: "6",
          left: "50%",
          position: "absolute",
          top: "50%",
          transform: "translate(-50%, -50%)",
          width: "6",
        },
      },
    },
  },
});
