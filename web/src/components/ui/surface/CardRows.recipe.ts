import { defineRecipe } from "@pandacss/dev";

export const cardRows = defineRecipe({
  className: "card-rows",
  base: {
    display: "flex",
    flexDirection: "column",
    alignItems: "start",
    width: "full",
    maxHeight: "min",
    gap: "4",
  },
});
