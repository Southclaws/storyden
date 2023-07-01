import { Box } from "@chakra-ui/react";
import { useCallback, useState } from "react";
import { createEditor } from "slate";
import { Editable, Slate, withReact } from "slate-react";

import { PostContent } from "src/api/openapi/schemas";

import { Leaf } from "../ContentComposer/render/Leaf";
import { Markdown } from "../Markdown";

type Props = {
  value: PostContent;
};

export function ContentViewer({ value }: Props) {
  const [editor] = useState(() => withReact(createEditor()));

  const renderLeaf = useCallback((props: any) => <Leaf {...props} />, []);

  switch (value.type) {
    case "text/markdown":
      return (
        <Box>
          <Markdown>{value.value ?? ""}</Markdown>
        </Box>
      );

    case "application/json":
      return (
        <Slate editor={editor} initialValue={JSON.parse(value.value)}>
          <Editable renderLeaf={renderLeaf} readOnly />
        </Slate>
      );

    default:
      console.error(`ContentViewer unexpected content type: ${value.type}`);
      return <></>;
  }
}
