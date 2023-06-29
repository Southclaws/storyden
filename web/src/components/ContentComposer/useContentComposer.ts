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
      return JSON.parse(props.initialValue);
    }

    return defaultValue;
  }, [props.initialValue]);

  function onChange(value: Descendant[]) {
    const isAstChange = editor.operations.some(
      (op) => "set_selection" !== op.type
    );

    if (isAstChange) {
      props.onChange(JSON.stringify(value));
    }
  }

  return { editor, initialValue, onChange };
}
