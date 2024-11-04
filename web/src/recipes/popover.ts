import { popoverAnatomy } from "@ark-ui/anatomy";
import { defineSlotRecipe } from "@pandacss/dev";

export const popover = defineSlotRecipe({
  className: "popover",
  slots: popoverAnatomy.keys(),
  base: {
    positioner: {
      position: "relative",
      maxHeight: "var(--available-height)",
      height: "min",
    },
    content: {
      maxHeight: "var(--available-height)",
      background: "bg.default",
      borderRadius: "l3",
      boxShadow: "lg",
      display: "flex",
      flexDirection: "column",
      maxWidth: "sm",
      zIndex: "popover",
      p: "4",
      _open: {
        animation: "fadeIn 0.25s ease-out",
      },
      _closed: {
        animation: "fadeOut 0.2s ease-out",
      },
      _hidden: {
        display: "none",
      },
    },
    title: {
      fontWeight: "medium",
      textStyle: "sm",
    },
    description: {
      overflowY: "scroll",
      color: "fg.muted",
      textStyle: "sm",
    },
    closeTrigger: {
      color: "fg.muted",
    },
    arrow: {
      "--arrow-size": "var(--sizes-3)",
      "--arrow-background": "var(--colors-bg-default)",
    },
    arrowTip: {
      borderTopWidth: "1px",
      borderLeftWidth: "1px",
    },
  },
});
