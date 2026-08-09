import { defineRecipe } from "@pandacss/dev";

export const multiSelectPicker = defineRecipe({
  className: "multi-select-picker",
  defaultVariants: {
    size: "sm",
  },
  variants: {
    size: {
      sm: {
        paddingInline: "0.5!",
      },
      md: {
        paddingInline: "1!",
      },
      lg: {
        paddingInline: "1.5!",
      },
    },
  },
});
