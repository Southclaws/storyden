import { defineRecipe } from "@pandacss/dev";

import { richCardShell } from "../rich-card/RichCard.recipe";

export const cardBox = defineRecipe({
  className: "card-box",
  jsx: ["CardBox"],
  base: {
    ...richCardShell,
    display: "flex",
    flexDirection: "column",
    gap: "1",
    width: "full",
    padding: "2",
  },
  defaultVariants: {
    kind: "default",
  },
  variants: {
    kind: {
      default: {},
      edge: {
        padding: "0",
      },
    },
  },
});
