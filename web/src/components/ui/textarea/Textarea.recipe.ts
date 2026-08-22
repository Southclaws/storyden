import { defineRecipe } from "@pandacss/dev";

import { inputBase, inputVariantVariants } from "../input/Input.recipe";

export const textarea = defineRecipe({
  className: "textarea",
  jsx: ["Textarea"],
  base: {
    ...inputBase,
    display: "block",
    lineHeight: "relaxed",
    resize: "vertical",
  },
  defaultVariants: {
    size: "md",
    variant: "outline",
  },
  variants: {
    size: {
      sm: { fontSize: "xs" },
      md: { fontSize: "sm" },
      lg: { fontSize: "md" },
    },
    variant: inputVariantVariants,
  },
  compoundVariants: [
    { size: "sm", variant: "outline", css: { px: "2", py: "1.5" } },
    { size: "md", variant: "outline", css: { px: "2.5", py: "2" } },
    { size: "lg", variant: "outline", css: { px: "3", py: "2.5" } },
  ],
});
