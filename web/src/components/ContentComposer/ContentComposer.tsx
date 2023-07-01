import { Box, BoxProps } from "@chakra-ui/react";
import { PropsWithChildren, useCallback } from "react";
import { Editable, Slate } from "slate-react";

import { Leaf } from "./render/Leaf";
import { Props, useContentComposer } from "./useContentComposer";

export function ContentComposer({
  disabled,
  children,
  ...props
}: PropsWithChildren<Props & Omit<BoxProps, "onChange">>) {
  const { editor, initialValue, onChange } = useContentComposer(props);

  const renderLeaf = useCallback((props: any) => <Leaf {...props} />, []);

  return (
    <Box id="rich-text-editor" w="full">
      <Slate editor={editor} initialValue={initialValue} onChange={onChange}>
        {children}
        <Editable
          renderLeaf={renderLeaf}
          readOnly={disabled}
          placeholder="Write your heart out..."
          style={{
            minHeight: props.minHeight ?? "24em",
            outline: "0px solid transparent",
            opacity: disabled ? 0.5 : 1,
          }}
        />
      </Slate>
    </Box>
  );
}
