import { Editor } from "slate";

import { Formats } from "../types";

export const isMarkActive = (editor: Editor, format: Formats) => {
  const marks = Editor.marks(editor);

  return marks ? marks[format] === true : false;
};

export function toggleMark(editor: Editor, format: Formats) {
  const isActive = isMarkActive(editor, format);

  if (isActive) {
    Editor.removeMark(editor, format);
    return false;
  } else {
    Editor.addMark(editor, format, true);
    return true;
  }
}
