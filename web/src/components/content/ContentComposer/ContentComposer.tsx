import { PropsWithChildren, useCallback } from "react";
import { Editor, Transforms } from "slate";
import { Editable, Slate } from "slate-react";

import { Box } from "@/styled-system/jsx";

import { FileDrop } from "./components/FileDrop/FileDrop";
import { Element } from "./render/Element";
import { Leaf } from "./render/Leaf";
import { Props, useContentComposer } from "./useContentComposer";
import { getURL } from "./utils";

export function ContentComposer({
  disabled,
  children,
  ...props
}: PropsWithChildren<Props>) {
  const { editor, initialValue, onChange } = useContentComposer(props);

  const renderLeaf = useCallback((props: any) => <Leaf {...props} />, []);
  const renderElement = useCallback((props: any) => <Element {...props} />, []);

  return (
    <Box
      id="rich-text-editor"
      className="typography"
      w="full"
      h="full"
      onDragOver={(e) => e.preventDefault()}
    >
      <Slate editor={editor} initialValue={initialValue} onChange={onChange}>
        {children}

        <FileDrop onComplete={props.onAssetUpload}>
          <Editable
            renderLeaf={renderLeaf}
            renderElement={renderElement}
            // onKeyDown={(event: React.KeyboardEvent<HTMLElement>) => {
            //   // NOTE: this hook prevents Slate from duplicating the previous
            //   // node (which results in images being duplicated.)

            //   if (event.key === "Enter") {
            //     event.preventDefault();

            //     console.log("enter", editor.selection);

            //     Transforms.insertNodes(editor, [
            //       {
            //         type: "paragraph",
            //         children: [
            //           {
            //             text: "",
            //           },
            //         ],
            //       },
            //     ]);
            //   }
            // }}
            readOnly={disabled}
            placeholder="Write your heart out..."
            style={{
              minHeight: props.minHeight ?? "8em",
              outline: "0px solid transparent",
              opacity: disabled ? 0.5 : 1,
            }}
          />
        </FileDrop>
      </Slate>
    </Box>
  );
}
