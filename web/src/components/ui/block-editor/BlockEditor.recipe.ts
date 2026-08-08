import { defineSlotRecipe } from "@pandacss/dev";

export const blockEditor = defineSlotRecipe({
  className: "block-editor",
  slots: [
    "root",
    "gutter",
    "handle",
    "menuTrigger",
    "menuAnchor",
    "menuContent",
    "content",
  ],
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
    menuTrigger: {
      width: "6",
      minWidth: "6",
      height: "6",
      padding: "0",
      color: "text.muted",
      cursor: "grab",
      _hover: {
        color: "text.default",
      },
      _expanded: {
        color: "text.default",
      },
      "&[data-dragging]": {
        cursor: "grabbing",
      },
    },
    menuAnchor: {
      position: "absolute",
      inset: "0",
      width: "full",
      height: "full",
      pointerEvents: "none",
    },
    menuContent: {
      minWidth: "40",
    },
    content: {
      width: "full",
      minWidth: "0",
    },
  },
});
