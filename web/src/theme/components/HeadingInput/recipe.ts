import { defineRecipe } from "@pandacss/dev";

export const headingInput = defineRecipe({
  className: "headingInput",
  base: {
    display: "inline-block",
    width: "full",
    fontSize: "3xl",
    overflowWrap: "break-word",
    wordBreak: "break-word",
    fontWeight: "semibold",
    cursor: "text",
    borderBottomColor: "border.default",
    borderBottomWidth: "1px",
    borderBottomStyle: "solid",
    _focus: {
      outline: "none",
      borderBottomColor: "accent.500",
      borderBottomWidth: "1px",
      borderBottomStyle: "solid",
    },
    _empty: {
      _before: {
        content: "attr(placeholder)",
        opacity: 0.3,
        color: "fg.default",
      },
    },
    _invalid: {
      borderBottomColor: "fg.destructive",
      borderBottomWidth: "1px",
      borderBottomStyle: "solid",
    },
  },
});
