import { BaseEditor, Element, Node, Path, Range, Transforms } from "slate";
import { ReactEditor } from "slate-react";

export function getURL(element: Element): string | undefined {
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

export const withExtensions = (editor: BaseEditor & ReactEditor) => {
  const { isVoid, normalizeNode, insertBreak, deleteBackward } = editor;

  editor.isVoid = (element) => {
    if (getURL(element)) {
      return true;
    }

    return isVoid(element);
  };

  editor.normalizeNode = ([node, path]) => {
    if (path.length === 0) {
      if (editor.children.length <= 1 && editor.string([0, 0]) === "") {
        Transforms.insertNodes(
          editor,
          {
            type: "paragraph",
            children: [{ text: "" }],
          },
          {
            at: path.concat(0),
            select: true,
          },
        );
      }
    }

    if (Element.isElement(node) && node.type === "paragraph") {
      for (const [child, childPath] of Node.children(editor, path)) {
        if (Element.isElement(child) && !editor.isInline(child)) {
          Transforms.unwrapNodes(editor, { at: childPath });
          return;
        }
      }
    }

    if (Element.isElement(node) && node.type === "paragraph") {
      if (editor.isVoid(node)) {
        Transforms.insertNodes(
          editor,
          {
            type: "paragraph",
            children: [{ text: "" }],
          },
          {},
        );
        return;
      }
    }

    return normalizeNode([node, path]);
  };

  editor.insertBreak = () => {
    if (!editor.selection || !Range.isCollapsed(editor.selection)) {
      return insertBreak();
    }

    const selectedNodePath = Path.parent(editor.selection.anchor.path);
    const selectedNode = Node.get(editor, selectedNodePath) as Element;

    if (editor.isVoid(selectedNode)) {
      const nextNodePath = Path.next(selectedNodePath);

      Transforms.insertNodes(
        editor,
        {
          type: "paragraph",
          children: [{ text: "" }],
        },
        {
          at: nextNodePath,
        },
      );

      Transforms.deselect(editor);
      Transforms.select(editor, nextNodePath);

      return;
    }

    insertBreak();
  };

  // if prev node is a void node, remove the current node and select the void node
  editor.deleteBackward = (unit) => {
    if (
      !editor.selection ||
      !Range.isCollapsed(editor.selection) ||
      editor.selection.anchor.offset !== 0
    ) {
      return deleteBackward(unit);
    }

    const parentPath = Path.parent(editor.selection.anchor.path);
    if (Path.hasPrevious(parentPath)) {
      const parentNode = Node.get(editor, parentPath);
      const parentIsEmpty = Node.string(parentNode).length === 0;

      if (parentIsEmpty && Path.hasPrevious(parentPath)) {
        const prevNodePath = Path.previous(parentPath);
        const prevNode = Node.get(editor, prevNodePath) as Element;
        if (editor.isVoid(prevNode)) {
          return Transforms.removeNodes(editor);
        }
      }
    }

    deleteBackward(unit);
  };

  return editor;
};
