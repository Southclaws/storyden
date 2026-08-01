import { popoverAnatomy } from "@ark-ui/react";
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
      background: "background.overlay",
      borderColor: "border.strong",
      borderRadius: "md",
      borderWidth: "thin",
      boxShadow: "overlay",
      color: "text.default",
      display: "flex",
      flexDirection: "column",
      maxWidth: "sm",
      zIndex: "popover",
      p: "2",
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
      color: "text.muted",
      textStyle: "sm",
    },
    closeTrigger: {
      color: "text.muted",
    },
    arrow: {
      "--arrow-size": "var(--sizes-3)",
      "--arrow-background": "var(--colors-background-overlay)",
    },
    arrowTip: {
      borderColor: "border.strong",
      borderTopWidth: "1px",
      borderLeftWidth: "1px",
    },
  },
});
