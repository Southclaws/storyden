import { Box, BoxProps } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { Editable, Slate } from "slate-react";

import { Props, useCompose } from "./useCompose";

export function Compose({
  children,
  ...props
}: PropsWithChildren<Props & BoxProps>) {
  const { editor, initialValue, onChange } = useCompose(props);

  return (
    <Box id="rich-text-editor" w="full">
      <Slate editor={editor} initialValue={initialValue} onChange={onChange}>
        <Editable
          placeholder="Write your heart out..."
          style={{
            minHeight: "24em",
            outline: "0px solid transparent",
          }}
        />
        {children}
      </Slate>
    </Box>
  );
}
