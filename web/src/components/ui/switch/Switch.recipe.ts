import { switchAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const switchRecipe = defineSlotRecipe({
  className: "switch",
  jsx: ["Switch", /Switch\.+/],
  slots: switchAnatomy.keys(),
  base: {
    root: {
      alignItems: "center",
      display: "flex",
      position: "relative",
      "&:has(input:focus-visible) [data-part=control]": {
        outline: "2px solid token(colors.accent.solid)",
        outlineOffset: "2px",
      },
    },
    control: {
      alignItems: "center",
      background: "background.controlTrack",
      borderRadius: "full",
      boxShadow: "inset 0 0 0 1px token(colors.border.default)",
      cursor: "pointer",
      display: "inline-flex",
      flexShrink: "0",
      position: "relative",
      transitionDuration: "normal",
      transitionProperty: "background",
      transitionTimingFunction: "default",
      _checked: {
        background: "accent.solid",
        boxShadow: "inset 0 0 0 1px token(colors.accent.solid)",
      },
      _disabled: {
        background: "background.controlDisabled",
        boxShadow: "inset 0 0 0 1px token(colors.border.disabled)",
        cursor: "not-allowed",
      },
    },
    label: {
      color: "text.default",
      fontWeight: "medium",
    },
    thumb: {
      background: "background.control",
      borderRadius: "full",
      boxShadow: "surface",
      transitionDuration: "normal",
      transitionProperty: "transform, background",
      transitionTimingFunction: "default",
      _checked: {
        transform: "translateX(100%)",
        background: "background.control",
      },
      _disabled: {
        background: "text.disabled",
      },
    },
  },
  defaultVariants: {
    size: "md",
  },
  variants: {
    size: {
      sm: {
        root: {
          gap: "2",
        },
        control: {
          width: "7",
          p: "0.5",
        },
        thumb: {
          width: "3",
          height: "3",
        },
        label: {
          fontSize: "sm",
          lineHeight: "1.25rem",
        },
      },
      md: {
        root: {
          gap: "3",
        },
        control: {
          width: "9",
          p: "0.5",
        },
        thumb: {
          width: "4",
          height: "4",
        },
        label: {
          fontSize: "md",
          lineHeight: "1.5rem",
        },
      },
      lg: {
        root: {
          gap: "4",
        },
        control: {
          width: "11",
          p: "0.5",
        },
        thumb: {
          width: "5",
          height: "5",
        },
        label: {
          fontSize: "lg",
          lineHeight: "1.75rem",
        },
      },
    },
  },
});
