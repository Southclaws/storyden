import {
  BaseEditor,
  Editor,
  Element,
  Node,
  Path,
  Range,
  Transforms,
} from "slate";
import { ReactEditor } from "slate-react";

import { CustomElement, ParagraphElement } from "./types";

export function getURL(element: CustomElement): string | undefined {
  if (element.children.length === 1 && element.children[0]) {
    const content = element.children[0];

    try {
      if (content.text.includes(" ")) {
        throw new Error();
      }

      const parsed = new URL(content.text);

      return parsed.toString();
    } catch (_) {
      return undefined;
    }
  }

  return undefined;
}

export const withCorrectVoidBehavior = (editor: BaseEditor & ReactEditor) => {
  const { deleteBackward, insertBreak } = editor;

  // if current selection is void node, insert a default node below
  editor.insertBreak = () => {
    console.log("insertBreak");

    if (!editor.selection || !Range.isCollapsed(editor.selection)) {
      return insertBreak();
    }

    const selectedNodePath = Path.parent(editor.selection.anchor.path);
    const selectedNode = Node.get(editor, selectedNodePath);
    if (Editor.isVoid(editor, selectedNode)) {
      Editor.insertNode(editor, {
        type: "paragraph",
        children: [{ text: "" }],
      });
      return;
    }

    insertBreak();
  };

  // if prev node is a void node, remove the current node if it's empty and select the void node
  editor.deleteBackward = (unit) => {
    console.log("deleteBackward");

    if (
      !editor.selection ||
      !Range.isCollapsed(editor.selection) ||
      editor.selection.anchor.offset !== 0
    ) {
      return deleteBackward(unit);
    }

    const parentPath = Path.parent(editor.selection.anchor.path);

    if (Path.hasPrevious(parentPath)) {
      const prevNodePath = Path.previous(parentPath);
      const prevNode = Node.get(editor, prevNodePath);
      if (Editor.isVoid(editor, prevNode)) {
        const parentNode = Node.get(editor, parentPath);
        const parentIsEmpty = Node.string(parentNode).length === 0;

        if (parentIsEmpty) {
          return Transforms.removeNodes(editor);
        } else {
          return Transforms.select(editor, prevNodePath);
        }
      }
    }

    deleteBackward(unit);
  };

  return editor;
};
