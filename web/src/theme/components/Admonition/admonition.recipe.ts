import { defineRecipe } from "@pandacss/dev";

export const admonition = defineRecipe({
  className: "admonition",
  base: {
    display: "flex",
    flexBasis: "0",
    flexGrow: "0",
    gap: "2",
    justifyContent: "space-between",
    borderColor: "blackAlpha.50",
    borderRadius: "lg",
    borderWidth: "1px",
    outline: 0,
    position: "relative",
    transitionDuration: "normal",
    transitionProperty: "box-shadow, border-color",
    transitionTimingFunction: "default",
    width: "full",
    padding: "2",
  },
  variants: {
    kind: {
      neutral: {
        backgroundColor: "whiteAlpha.600",
      },
      success: {
        backgroundColor: "green.100",
        color: "accent.text.500",
      },
      failure: {
        backgroundColor: "bg.destructive",
        color: "fg.default",
      },
    },
  },
});
