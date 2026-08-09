import { defineRecipe } from "@pandacss/dev";

export const text = defineRecipe({
  className: "text",
  jsx: ["Text"],
  variants: {
    variant: {
      body: {
        color: "text.default",
        textStyle: "body",
      },
      supporting: {
        color: "text.muted",
        textStyle: "supporting",
      },
      metadata: {
        color: "text.muted",
        textStyle: "metadata",
      },
    },
  },
  defaultVariants: {
    variant: "body",
  },
});
