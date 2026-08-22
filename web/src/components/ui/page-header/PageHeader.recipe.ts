import { defineSlotRecipe } from "@pandacss/dev";

export const pageHeader = defineSlotRecipe({
  className: "page-header",
  slots: [
    "root",
    "navigation",
    "row",
    "heading",
    "titleGroup",
    "titleRow",
    "actions",
  ],
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
      flexDirection: { base: "column", md: "row" },
    },
    heading: {
      display: "flex",
      width: "full",
      minWidth: "0",
      alignItems: "start",
      gap: "1",
    },
    titleGroup: {
      display: "flex",
      minWidth: "0",
      flexDirection: "column",
      gap: "1",
      "& > .text": {
        overflowWrap: "anywhere",
      },
    },
    titleRow: {
      display: "flex",
      minWidth: "0",
      alignItems: "center",
      flexWrap: "wrap",
      gap: "2",
    },
    actions: {
      display: "flex",
      flex: "none",
      alignItems: "center",
      alignSelf: { base: "end", md: "start" },
      flexWrap: "wrap",
      justifyContent: "end",
      gap: "1",
    },
  },
});
