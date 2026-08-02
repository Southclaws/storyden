import { defineSlotRecipe } from "@pandacss/dev";

export const pageHeader = defineSlotRecipe({
  className: "page-header",
  slots: ["root", "navigation", "row", "heading", "titleGroup", "actions"],
  base: {
    root: {
      display: "flex",
      width: "full",
      minWidth: "0",
      flexDirection: "column",
      gap: "1",
    },
    navigation: {
      display: "flex",
      width: "full",
      minWidth: "0",
      alignItems: "center",
      gap: "1",
    },
    row: {
      display: "flex",
      width: "full",
      minWidth: "0",
      alignItems: "start",
      justifyContent: "space-between",
      gap: "2",
    },
    heading: {
      display: "flex",
      minWidth: "0",
      alignItems: "start",
      gap: "1",
    },
    titleGroup: {
      display: "flex",
      minWidth: "0",
      flexDirection: "column",
      gap: "1",
    },
    actions: {
      display: "flex",
      flex: "none",
      alignItems: "center",
      gap: "1",
    },
  },
});
