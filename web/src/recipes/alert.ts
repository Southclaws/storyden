import { defineSlotRecipe } from "@pandacss/dev";

export const alert = defineSlotRecipe({
  className: "alert",
  slots: ["root", "content", "description", "icon", "title"],
  base: {
    root: {
      background: "bg.muted",
      borderColor: "border.muted",
      borderWidth: "thin",
      borderRadius: "l3",
      display: "flex",
      gap: "2",
      p: "2",
      width: "full",
    },
    content: {
      display: "flex",
      flexDirection: "column",
      gap: "1",
    },
    description: {
      color: "fg.muted",
      textStyle: "sm",
    },
    icon: {
      color: "fg.warning",
      flexShrink: "0",
      pt: "0.5",
      width: "4",
      height: "4",
    },
    title: {
      color: "fg.warning",
      fontWeight: "semibold",
      textStyle: "sm",
    },
  },
});
