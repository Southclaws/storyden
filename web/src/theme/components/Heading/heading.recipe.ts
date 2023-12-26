import { defineRecipe } from "@pandacss/dev";

export const heading = defineRecipe({
  className: "heading",
  base: {
    fontWeight: "semibold",
  },
  defaultVariants: {
    size: "md",
  },
  variants: {
    size: {
      xs: {
        fontSize: "heading.6",
      },
      sm: {
        fontSize: "heading.5",
      },
      md: {
        fontSize: "heading.4",
      },
      lg: {
        fontSize: "heading.3",
      },
      xl: {
        fontSize: "heading.2",
      },
      "2xl": {
        fontSize: "heading.1",
      },
    },
  },
});
