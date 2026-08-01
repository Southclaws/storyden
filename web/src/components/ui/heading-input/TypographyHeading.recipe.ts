import { defineRecipe } from "@pandacss/dev";

export const typographyHeading = defineRecipe({
  className: "typography-heading",
  base: {
    fontWeight: "semibold",
  },
  defaultVariants: {
    size: "md",
  },

  variants: {
    size: {
      md: {
        textStyle: "xl",
      },
      lg: {
        textStyle: "2xl",
      },
    },
  },
});
