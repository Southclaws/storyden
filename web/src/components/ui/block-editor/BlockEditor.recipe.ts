import { defineSlotRecipe } from "@pandacss/dev";

export const blockEditor = defineSlotRecipe({
  className: "block-editor",
  slots: ["root", "gutter", "handle", "content"],
  base: {
    root: {
      position: "relative",
      display: "flex",
      width: "full",
      minWidth: "0",
    },
    gutter: {
      position: "absolute",
      insetBlockStart: "0",
      insetInlineStart: "-7",
      display: "none",
      width: "6",
      height: "full",
      alignItems: "start",
      padding: "0",
      md: {
        display: "flex",
      },
    },
    handle: {
      position: "relative",
      display: "flex",
      width: "full",
      minHeight: "6",
      alignItems: "start",
      justifyContent: "center",
      borderRadius: "sm",
      color: "text.muted",
      cursor: "grab",
      opacity: "full",
    },
    content: {
      width: "full",
      minWidth: "0",
    },
  },
});
