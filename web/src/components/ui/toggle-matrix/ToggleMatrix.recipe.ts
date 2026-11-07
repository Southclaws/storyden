import { defineSlotRecipe } from "@pandacss/dev";

export const toggleMatrix = defineSlotRecipe({
  className: "toggleMatrix",
  slots: ["root", "columnHeader", "rowHeader", "description", "cell", "item"],
  base: {
    root: {
      tableLayout: "fixed",
    },
    columnHeader: {
      width: { base: "16", sm: "20" },
      textAlign: "center",
      color: "text.muted",
    },
    rowHeader: {
      color: "text.default",
      fontWeight: "medium",
      overflowWrap: "anywhere",
    },
    description: {
      color: "text.muted",
      fontWeight: "normal",
      marginTop: "1",
    },
    cell: {
      textAlign: "center",
    },
    item: {
      width: "full",
      borderWidth: "thin",
      borderColor: "border.default",
    },
  },
});
