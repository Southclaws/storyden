import { defineRecipe } from "@pandacss/dev";

import { button } from "../Button/button.recipe";

export const link = defineRecipe({
  className: "link",
  base: {
    ...button.base,
  },
  defaultVariants: {
    kind: "neutral",
    size: "md",
  },
  variants: {
    ...button.variants,
  },
} as any);
