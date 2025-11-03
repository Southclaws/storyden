import { collapsibleAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const collapsible = defineSlotRecipe({
  className: "collapsible",
  slots: collapsibleAnatomy.keys(),
  base: {
    root: {
      alignItems: "flex-start",
      display: "flex",
      flexDirection: "column",
      gap: "1",
      width: "full",
    },
    content: {
      overflow: "hidden",
      width: "full",
      _open: {
        animation: "collapse-in",
      },
      _closed: {
        animation: "collapse-out",
      },
    },
  },
});
