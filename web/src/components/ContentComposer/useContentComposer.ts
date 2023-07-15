import { useMemo, useRef } from "react";
import {
  BaseEditor,
  Descendant,
  Editor,
  Transforms,
  createEditor,
} from "slate";
import { ReactEditor, withReact } from "slate-react";

import { Asset } from "src/api/openapi/schemas";

import { deserialise, serialise } from "./serialisation";

export type Props = {
  resetKey?: string;
  disabled?: boolean;
  minHeight?: string;
  initialValue?: string;
  onChange: (value: string) => void;
  onAssetUpload: (asset: Asset) => void;
};

const defaultValue: Descendant[] = [
  {
    type: "paragraph",
    children: [{ text: "" }],
  },
];

export function useContentComposer(props: Props) {
  const editorRef = useRef<BaseEditor & ReactEditor>();
  if (!editorRef.current) {
    editorRef.current = withReact(createEditor());
  }
  const editor = editorRef.current;

  const initialValue: Descendant[] = useMemo(() => {
    if (props.initialValue) {
      return deserialise(props.initialValue);
    }

    return defaultValue;
  }, [props.initialValue]);

  useMemo(() => {
    if (!props.resetKey) return;

    Transforms.delete(editor, {
      at: {
        anchor: Editor.start(editor, []),
        focus: Editor.end(editor, []),
      },
    });

    // Disable this error because, despite not using the "resetKey" prop, we're
    // using the behaviour of useMemo to clear the input when the value changes.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [props.resetKey, editor]);

  function onChange(value: Descendant[]) {
    const isAstChange = editor.operations.some(
      (op) => "set_selection" !== op.type
    );

    if (isAstChange) {
      props.onChange(serialise(value));
    }
  }

  return { editor, initialValue, onChange };
}
