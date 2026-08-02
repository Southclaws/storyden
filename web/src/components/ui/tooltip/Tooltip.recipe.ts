import { tooltipAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const tooltip = defineSlotRecipe({
  className: "tooltip",
  slots: tooltipAnatomy.keys(),
  base: {
    content: {
      background: "background.overlay",
      borderColor: "border.strong",
      borderRadius: "sm",
      borderWidth: "thin",
      boxShadow: "floating",
      color: "text.default",
      fontWeight: "semibold",
      px: "3",
      py: "2",
      fontSize: "xs",
      lineHeight: "1.125rem",
      maxWidth: "2xs",
      zIndex: "tooltip",
      _open: {
        animation: "fadeIn 0.25s ease-out",
      },
      _closed: {
        animation: "fadeOut 0.2s ease-out",
      },
    },
  },
});
