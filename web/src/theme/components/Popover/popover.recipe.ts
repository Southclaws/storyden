import { popoverAnatomy } from "@ark-ui/anatomy";
import { defineSlotRecipe } from "@pandacss/dev";

export const popover = defineSlotRecipe({
  className: "popover",
  slots: popoverAnatomy.keys(),
  base: {
    positioner: {
      position: "relative",
    },
    content: {
      zIndex: "popover",
      borderRadius: "lg",
      boxShadow: "lg",
      display: "flex",
      flexDirection: "column",
      maxWidth: "sm",
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
      color: "fg.muted",
      textStyle: "sm",
    },
    closeTrigger: {
      color: "fg.muted",
    },
  },
});
