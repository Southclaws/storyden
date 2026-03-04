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
      backgroundColor: "bg.muted",
    },
    range: {
      height: "full",
      backgroundColor: "colorPalette.default",
      borderRadius: "full",
      transition: "all",
    },
    label: {
      color: "fg.default",
      fontWeight: "medium",
    },
    valueText: {
      color: "fg.muted",
      textAlign: "center",
    },
    circle: {
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
    },
    circleTrack: {
      stroke: "bg.muted",
      strokeWidth: "4px",
    },
    circleRange: {
      stroke: "colorPalette.default",
      strokeWidth: "4px",
      transitionProperty: "stroke-dasharray, stroke",
      transitionDuration: "0.6s",
    },
  },
  defaultVariants: {
    size: "md",
  },
  variants: {
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
          textStyle: "xs",
        },
        valueText: {
          textStyle: "xs",
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
          textStyle: "sm",
        },
        valueText: {
          textStyle: "sm",
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
          textStyle: "md",
        },
        valueText: {
          textStyle: "sm",
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
