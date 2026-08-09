import { progressAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const progress = defineSlotRecipe({
  className: "progress",
  slots: progressAnatomy.keys(),
  jsx: ["ProgressCircle", "ProgressHorizontal"],
  base: {
    root: {
      display: "flex",
      flexDirection: "column",
      gap: "1.5",
      colorPalette: "accent",
    },
    track: {
      borderRadius: "full",
      overflow: "hidden",
      backgroundColor: "background.controlTrack",
    },
    range: {
      height: "full",
      backgroundColor: "colorPalette.solid",
      borderRadius: "full",
      transition: "all",
    },
    label: {
      color: "text.default",
      fontWeight: "medium",
    },
    valueText: {
      color: "text.subtle",
      textAlign: "center",
    },
    circle: {
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
    },
    circleTrack: {
      stroke: "background.controlTrack",
      strokeWidth: "4px",
    },
    circleRange: {
      stroke: "colorPalette.solid",
      strokeWidth: "4px",
      transitionProperty: "stroke-dasharray, stroke",
      transitionDuration: "0.6s",
    },
  },
  defaultVariants: {
    shape: "horizontal",
    size: "md",
  },
  variants: {
    shape: {
      circle: {
        root: {
          alignItems: "center",
          width: "fit",
        },
      },
      horizontal: {
        root: {
          width: "full",
        },
      },
    },
    size: {
      sm: {
        root: {
          "--size": "var(--sizes-5)",
          "--thickness": "var(--border-widths-thick)",
        },
        track: {
          height: "1",
        },
        label: {
          fontSize: "xs",
          lineHeight: "1.125rem",
        },
        valueText: {
          fontSize: "xs",
          lineHeight: "1.125rem",
        },
        circle: {
          width: "9",
          height: "9",

          "& svg": {
            width: "9",
            height: "9",
          },
        },
        circleTrack: {
          strokeWidth: "3px",
        },
        circleRange: {
          strokeWidth: "3px",
        },
      },
      md: {
        root: {
          "--size": "var(--sizes-8)",
          "--thickness": "var(--sizes-1)",
        },
        track: {
          height: "2",
        },
        label: {
          fontSize: "sm",
          lineHeight: "1.25rem",
        },
        valueText: {
          fontSize: "sm",
          lineHeight: "1.25rem",
          // optical alignment because of % character
          paddingLeft: "0.5",
        },
        circle: {
          width: "12",
          height: "12",
          "& svg": {
            width: "12",
            height: "12",
          },
        },
        circleTrack: {
          strokeWidth: "4px",
        },
        circleRange: {
          strokeWidth: "4px",
        },
      },
      lg: {
        root: {
          "--size": "var(--sizes-12)",
          "--thickness": "var(--sizes-2)",
        },

        track: {
          height: "3",
        },
        label: {
          fontSize: "md",
          lineHeight: "1.5rem",
        },
        valueText: {
          fontSize: "sm",
          lineHeight: "1.25rem",
          // optical alignment because of % character
          paddingLeft: "2",
        },
        circle: {
          width: "16",
          height: "16",
          "& svg": {
            width: "16",
            height: "16",
          },
        },
        circleTrack: {
          strokeWidth: "5px",
        },
        circleRange: {
          strokeWidth: "5px",
        },
      },
    },
  },
});
