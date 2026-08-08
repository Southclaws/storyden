import { defineRecipe } from "@pandacss/dev";

export const admonition = defineRecipe({
  className: "admonition",
  base: {
    display: "flex",
    flexBasis: "0",
    flexGrow: "0",
    gap: "2",
    justifyContent: "space-between",
    backgroundColor: "background.inset",
    borderColor: "border.default",
    borderRadius: "lg",
    borderWidth: "1px",
    color: "text.default",
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
      neutral: {},
      success: {
        backgroundColor: "status.success.surface",
        borderColor: "status.success.border",
        color: "status.success.content",
      },
      failure: {
        backgroundColor: "status.danger.surface",
        borderColor: "status.danger.border",
        color: "status.danger.content",
      },
    },
  },
});
