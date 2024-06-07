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
      xs: {
        fontSize: "heading.6",
        "--font-size": "fontSizes.heading.6",
        "--font-size-fluid-scale": "4cqi",
        // NOTE: Does not scale at all.
        "--font-size-fluid-min": "fontSizes.heading.6",
      },
      sm: {
        fontSize: "heading.5",
        "--font-size": "fontSizes.heading.5",
        "--font-size-fluid-scale": "4cqi",
        "--font-size-fluid-min": "fontSizes.heading.6",
      },
      md: {
        fontSize: "heading.4",
        "--font-size": "fontSizes.heading.4",
        "--font-size-fluid-scale": "4cqi",
        "--font-size-fluid-min": "fontSizes.heading.5",
      },
      lg: {
        fontSize: "heading.3",
        "--font-size": "fontSizes.heading.3",
        "--font-size-fluid-scale": "4cqi",
        "--font-size-fluid-min": "fontSizes.heading.4",
      },
      xl: {
        fontSize: "heading.2",
        "--font-size": "fontSizes.heading.2",
        "--font-size-fluid-scale": "4cqi",
        "--font-size-fluid-min": "fontSizes.heading.3",
      },
      "2xl": {
        fontSize: "heading.1",
        "--font-size": "fontSizes.heading.1",
        "--font-size-fluid-scale": "4cqi",
        "--font-size-fluid-min": "fontSizes.heading.2",
      },
    },
  },
});
