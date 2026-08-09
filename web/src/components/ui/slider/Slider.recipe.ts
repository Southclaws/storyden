import { sliderAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const slider = defineSlotRecipe({
  className: "slider",
  slots: sliderAnatomy.keys(),
  base: {
    root: {
      colorPalette: "accent",
      display: "flex",
      flexDirection: "column",
      gap: "1",
      width: "full",
    },
    control: {
      position: "relative",
      display: "flex",
      alignItems: "center",
    },
    track: {
      backgroundColor: "background.controlTrack",
      borderRadius: "full",
      overflow: "hidden",
      flex: "1",
    },
    range: {
      background: "accent.default",
    },
    thumb: {
      background: "background.control",
      borderColor: "accent.default",
      borderRadius: "full",
      borderWidth: "2px",
      boxShadow: "surface",
      outline: "none",
      zIndex: "1",
      transition: "box-shadow 0.2s, transform 0.2s",
      _focusVisible: {
        boxShadow: "0 0 0 3px token(colors.accent.4)",
        outline: "2px solid token(colors.accent.default)",
        outlineOffset: "2px",
      },
      _hover: {
        transform: "scale(1.1)",
      },
      _active: {
        transform: "scale(0.95)",
      },
    },
    label: {
      color: "text.default",
      fontWeight: "medium",
    },
    markerGroup: {
      mt: "-1",
      minHeight: "1lh",
    },
    marker: {
      "--before-background": {
        _osLight: "white",
        _osDark: "black",
      },
      color: "text.subtle",
      _before: {
        background: "white",
        borderRadius: "full",
        content: "''",
        display: "block",
        left: "50%",
        position: "relative",
        transform: "translateX(-50%)",
      },
      _underValue: {
        _before: {
          background: "var(--before-background)",
        },
      },
    },
  },
  defaultVariants: {
    size: "md",
  },
  variants: {
    size: {
      sm: {
        control: {
          height: "4",
        },
        range: {
          height: "1.5",
        },
        track: {
          height: "1.5",
        },
        thumb: {
          height: "4",
          width: "4",
        },
        marker: {
          _before: {
            height: "1",
            top: "-2.5",
            width: "1",
          },
          fontSize: "xs",
          lineHeight: "1.125rem",
        },
        label: {
          fontSize: "sm",
          lineHeight: "1.25rem",
        },
      },
      md: {
        control: {
          height: "5",
        },
        range: {
          height: "2",
        },
        track: {
          height: "2",
        },
        thumb: {
          height: "5",
          width: "5",
        },
        marker: {
          _before: {
            height: "1",
            top: "-3",
            width: "1",
          },
          fontSize: "xs",
          lineHeight: "1.125rem",
        },
        label: {
          fontSize: "sm",
          lineHeight: "1.25rem",
        },
      },
      lg: {
        control: {
          height: "6",
        },
        range: {
          height: "2.5",
        },
        track: {
          height: "2.5",
        },
        thumb: {
          height: "6",
          width: "6",
        },
        marker: {
          _before: {
            height: "1.5",
            top: "-15px",
            width: "1.5",
          },
          fontSize: "sm",
          lineHeight: "1.25rem",
        },
        label: {
          fontSize: "md",
          lineHeight: "1.5rem",
        },
      },
    },
  },
});
