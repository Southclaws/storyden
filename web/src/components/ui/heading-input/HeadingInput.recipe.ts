import { defineRecipe } from "@pandacss/dev";

export const headingInput = defineRecipe({
  className: "headingInput",
  base: {
    color: "text.default",
    display: "inline-block",
    fontSize: "xl",
    width: "full",
    overflowWrap: "break-word",
    wordBreak: "break-word",
    fontWeight: "semibold",
    lineHeight: "normal",
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
        color: "text.default",
      },
    },
    _invalid: {
      borderBottomColor: "status.danger.content",
      borderBottomWidth: "1px",
      borderBottomStyle: "solid",
    },
  },
});
