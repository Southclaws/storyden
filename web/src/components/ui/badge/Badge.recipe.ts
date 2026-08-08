import { defineRecipe } from "@pandacss/dev";

export const badgeColorPalettes = [
  "accent",
  "gray",
  "blue",
  "green",
  "amber",
  "orange",
  "pink",
  "red",
  "slate",
  "tomato",
] as const;

export const badge = defineRecipe({
  className: "badge",
  jsx: ["Badge"],
  staticCss: [{ size: ["sm", "md", "lg"] }],
  base: {
    alignItems: "center",
    borderRadius: "full",
    display: "inline-flex",
    fontWeight: "medium",
    userSelect: "none",
    whiteSpace: "nowrap",
  },
  defaultVariants: {
    variant: "subtle",
    size: "md",
  },
  variants: {
    variant: {
      solid: {
        background: "colorPalette.default",
        color: "colorPalette.fg",
      },
      subtle: {
        background: "colorPalette.3",
        borderColor: "colorPalette.6",
        borderWidth: "1px",
        color: "colorPalette.12",
        "& svg": {
          color: "colorPalette.12",
        },
      },
      outline: {
        color: "colorPalette.12",
        borderWidth: "2px",
        borderColor: "colorPalette.7",
      },
    },
    size: {
      sm: {
        fontSize: "xs",
        lineHeight: "1.125rem",
        px: "2",
        h: "5",
        gap: "1",
        "& svg": {
          width: "3",
          height: "3",
        },
      },
      md: {
        fontSize: "xs",
        lineHeight: "1.125rem",
        px: "2.5",
        h: "6",
        gap: "1.5",
        "& svg": {
          width: "4",
          height: "4",
        },
      },
      lg: {
        fontSize: "sm",
        lineHeight: "1.25rem",
        px: "3",
        h: "7",
        gap: "1.5",
        "& svg": {
          width: "4",
          height: "4",
        },
      },
    },
  },
});
