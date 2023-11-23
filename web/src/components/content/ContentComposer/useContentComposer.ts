"use client";

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

import { withEmbeds } from "./render/Element";
import { deserialise, serialise } from "./serialisation";
import { withCorrectVoidBehavior } from "./utils";

export type Props = {
  resetKey?: string;
  disabled?: boolean;
  minHeight?: string;
  initialValue?: string;
  onChange: (value: string) => void;
  onAssetUpload?: (asset: Asset) => void;
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
    editorRef.current = withCorrectVoidBehavior(
      withEmbeds(withReact(createEditor())),
    );
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
  }, [props.resetKey, editor]);

  function onChange(value: Descendant[]) {
    const isAstChange = editor.operations.some(
      (op) => "set_selection" !== op.type,
    );

    if (isAstChange) {
      props.onChange(serialise(value));
    }
  }

  return { editor, initialValue, onChange };
}
