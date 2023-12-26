import { defineRecipe } from "@pandacss/dev";

export const titleInput = defineRecipe({
  className: "titleInput",
  base: {
    display: "inline-block",
    contentEditable: true,
    //
    // NOTE: We're doing a bit of a hack here in order to make this
    // field look nice and behave like the Substack title editor.
    //
    // More info:
    //
    // https://medium.com/programming-essentials/good-to-know-about-the-state-management-of-a-contenteditable-element-in-react-adb4f933df12
    //
    suppressContentEditableWarning: true,
    width: "full",
    fontSize: "3xl",
    overflowWrap: "break-word",
    wordBreak: "break-word",
    fontWeight: "semibold",
    placeholder: "Title",
    cursor: "text",
    _focus: {
      outline: "none",
    },
    _empty: {
      _before: {
        content: "attr(placeholder)",
        opacity: 0.3,
        color: "fg.default",
      },
    },
  },
});
