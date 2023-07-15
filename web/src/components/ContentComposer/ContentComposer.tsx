import { Box, BoxProps } from "@chakra-ui/react";
import { PropsWithChildren, useCallback } from "react";
import { Transforms } from "slate";
import { Editable, Slate } from "slate-react";

import { FileDrop } from "./components/FileDrop/FileDrop";
import { Element } from "./render/Element";
import { Leaf } from "./render/Leaf";
import { Props, useContentComposer } from "./useContentComposer";

export function ContentComposer({
  disabled,
  children,
  ...props
}: PropsWithChildren<Props & Omit<BoxProps, "onChange">>) {
  const { editor, initialValue, onChange } = useContentComposer(props);

  const renderLeaf = useCallback((props: any) => <Leaf {...props} />, []);
  const renderElement = useCallback((props: any) => <Element {...props} />, []);

  return (
    <Box id="rich-text-editor" w="full" onDragOver={(e) => e.preventDefault()}>
      <Slate editor={editor} initialValue={initialValue} onChange={onChange}>
        {children}

        <FileDrop onComplete={props.onAssetUpload}>
          <Editable
            renderLeaf={renderLeaf}
            renderElement={renderElement}
            onKeyDown={(event: React.KeyboardEvent<HTMLElement>) => {
              // NOTE: this hook prevents Slate from duplicating the previous
              // node (which results in images being duplicated.)

              if (event.key === "Enter") {
                event.preventDefault();

                Transforms.insertNodes(editor, [
                  {
                    type: "paragraph",
                    children: [
                      {
                        text: "",
                      },
                    ],
                  },
                ]);
              }
            }}
            readOnly={disabled}
            placeholder="Write your heart out..."
            style={{
              minHeight: props.minHeight ?? "24em",
              outline: "0px solid transparent",
              opacity: disabled ? 0.5 : 1,
            }}
          />
        </FileDrop>
      </Slate>
    </Box>
  );
}
