import { useMemo, useState } from "react";
import { Descendant, Editor, Transforms, createEditor } from "slate";
import { withReact } from "slate-react";

import { deserialise, serialise } from "./serialisation";

export type Props = {
  resetKey?: string;
  disabled?: boolean;
  minHeight?: string;
  initialValue?: string;
  onChange: (value: string) => void;
};

const defaultValue: Descendant[] = [
  {
    type: "paragraph",
    children: [{ text: "" }],
  },
];

export function useContentComposer(props: Props) {
  const [editor] = useState(() => withReact(createEditor()));

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
