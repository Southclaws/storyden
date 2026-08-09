import { defineRecipe } from "@pandacss/dev";

const withoutIntent =
  "&:not(.button--intent_success):not(.button--intent_warning):not(.button--intent_destructive)";

export const button = defineRecipe({
  className: "button",
  jsx: ["Button", "IconButton", "SubmitButton"],
  base: {
    alignItems: "center",
    appearance: "none",
    borderRadius: "sm",
    colorPalette: "gray",
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
        [withoutIntent]: {
          colorPalette: "accent",
        },
        _hover: {
          background: "colorPalette.default",
          color: "colorPalette.fg",
          outline: "1px solid",
          outlineColor: "colorPalette.fg/20",
          outlineOffset: "-1px",
        },
        _focusVisible: {
          outline: "2px solid",
          outlineColor: "colorPalette.text",
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
        background: "background.control",
        borderWidth: "1px",
        borderColor: "colorPalette.7",
        color: "colorPalette.text",
        _hover: {
          background: "colorPalette.3",
          borderColor: "colorPalette.8",
          color: "colorPalette.text",
        },
        _disabled: {
          background: "background.controlDisabled",
          borderColor: "border.disabled",
          color: "text.disabled",
          cursor: "not-allowed",
          _hover: {
            background: "background.controlDisabled",
            borderColor: "border.disabled",
            color: "text.disabled",
          },
        },
        _focusVisible: {
          outline: "2px solid",
          outlineColor: "colorPalette.text",
          outlineOffset: "2px",
        },
        _selected: {
          background: "colorPalette.default",
          borderColor: "colorPalette.default",
          color: "colorPalette.fg",
          _hover: {
            background: "colorPalette.default",
            borderColor: "colorPalette.default",
          },
        },
      },
      ghost: {
        color: "text.muted",
        [withoutIntent]: {
          colorPalette: "accent",
        },
        _hover: {
          background: "colorPalette.3",
          color: "colorPalette.text",
        },
        _selected: {
          background: "colorPalette.4",
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
          outlineColor: "colorPalette.text",
          outlineOffset: "2px",
        },
      },
      subtle: {
        background: "colorPalette.3/80",
        color: "colorPalette.text",
        backdropBlur: "sm",
        backdropFilter: "auto",
        _hover: {
          background: "colorPalette.4",
          color: "colorPalette.text",
        },
        _focusVisible: {
          outline: "2px solid",
          outlineColor: "colorPalette.text",
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
    intent: {
      success: {
        colorPalette: "green",
      },
      warning: {
        colorPalette: "amber",
      },
      destructive: {
        colorPalette: "red",
      },
    },
  },
});
