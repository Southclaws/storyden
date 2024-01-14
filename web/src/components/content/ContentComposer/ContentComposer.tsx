import { PropsWithChildren, useCallback } from "react";
import { Editable, Slate } from "slate-react";

import { FileDrop } from "../FileDrop/FileDrop";

import { Box } from "@/styled-system/jsx";

import { Element } from "./render/Element";
import { Leaf } from "./render/Leaf";
import { Props, useContentComposer } from "./useContentComposer";

export function ContentComposer({
  disabled,
  children,
  ...props
}: PropsWithChildren<Props>) {
  const { editor, initialValue, onChange, handleAssetUpload } =
    useContentComposer(props);

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

        <FileDrop onComplete={handleAssetUpload}>
          <Editable
            renderLeaf={renderLeaf}
            renderElement={renderElement}
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
