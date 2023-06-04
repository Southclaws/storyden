import { useMemo, useState } from "react";
import { BaseEditor, Descendant, createEditor } from "slate";
import { ReactEditor, withReact } from "slate-react";

type CustomElement = { type: "paragraph"; children: CustomText[] };
type CustomText = { text: string };

declare module "slate" {
  interface CustomTypes {
    Editor: BaseEditor & ReactEditor;
    Element: CustomElement;
    Text: CustomText;
  }
}

export type Props = {
  initialValue?: string;
  onChange: (value: unknown) => void;
};

const defaultValue: Descendant[] = [
  {
    type: "paragraph",
    children: [{ text: "Write your heart out..." }],
  },
];

export function useCompose(props: Props) {
  const [editor] = useState(() => withReact(createEditor()));

  const initialValue = useMemo(
    () => (props.initialValue ? JSON.parse(props.initialValue) : defaultValue),
    [props.initialValue]
  );

  function onChange(value: unknown) {
    const isAstChange = editor.operations.some(
      (op) => "set_selection" !== op.type
    );

    if (isAstChange) {
      const content = JSON.stringify(value);

      props.onChange(content);
    }
  }

  return { editor, initialValue, onChange };
}
