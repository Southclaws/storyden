"use client";

import { useCallback, useState } from "react";
import { createEditor } from "slate";
import { Editable, Slate, withReact } from "slate-react";

import { PostContent } from "src/api/openapi/schemas";

import { Element } from "../ContentComposer/render/Element";
import { Leaf } from "../ContentComposer/render/Leaf";
import { deserialise } from "../ContentComposer/serialisation";

type Props = {
  value: PostContent;
};

export function ContentViewer({ value }: Props) {
  const [editor] = useState(() => withReact(createEditor()));

  const renderLeaf = useCallback((props: any) => <Leaf {...props} />, []);
  const renderElement = useCallback((props: any) => <Element {...props} />, []);

  return (
    <Slate editor={editor} initialValue={deserialise(value)}>
      <Editable
        className="typography content-viewer"
        renderLeaf={renderLeaf}
        renderElement={renderElement}
        readOnly
      />
    </Slate>
  );
}
