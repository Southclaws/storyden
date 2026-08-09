import { defineSlotRecipe } from "@pandacss/dev";

export const alert = defineSlotRecipe({
  className: "alert",
  slots: ["root", "content", "description", "icon", "title"],
  base: {
    root: {
      background: "background.inset",
      borderColor: "border.default",
      borderWidth: "thin",
      borderRadius: "md",
      color: "text.default",
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
      color: "text.muted",
      fontSize: "sm",
      lineHeight: "1.25rem",
    },
    icon: {
      color: "status.warning.content",
      flexShrink: "0",
      pt: "0.5",
      width: "4",
      height: "4",
    },
    title: {
      color: "status.warning.content",
      fontWeight: "semibold",
      fontSize: "sm",
      lineHeight: "1.25rem",
    },
  },
});
