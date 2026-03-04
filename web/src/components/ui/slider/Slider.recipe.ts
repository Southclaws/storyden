import { sliderAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const slider = defineSlotRecipe({
  className: "slider",
  slots: sliderAnatomy.keys(),
  base: {
    root: {
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
      backgroundColor: "bg.emphasized",
      borderRadius: "full",
      overflow: "hidden",
      flex: "1",
    },
    range: {
      background: "colorPalette.default",
    },
    thumb: {
      background: "bg.default",
      borderColor: "colorPalette.default",
      borderRadius: "full",
      borderWidth: "2px",
      boxShadow: "sm",
      outline: "none",
      zIndex: "1",
      transition: "box-shadow 0.2s, transform 0.2s",
      _focusVisible: {
        boxShadow: "0 0 0 3px token(colors.colorPalette.a4)",
        outline: "2px solid token(colors.colorPalette.default)",
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
      color: "fg.default",
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
      color: "fg.muted",
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
          textStyle: "xs",
        },
        label: {
          textStyle: "sm",
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
          textStyle: "xs",
        },
        label: {
          textStyle: "sm",
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
          textStyle: "sm",
        },
        label: {
          textStyle: "md",
        },
      },
    },
  },
});
