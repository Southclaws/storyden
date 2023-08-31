"use client";

import { Box } from "@chakra-ui/react";
import { PropsWithChildren } from "react";

import { Props, useFileDrop } from "./useFileDrop";

export function FileDrop({ children, ...props }: PropsWithChildren<Props>) {
  const { onDragStart, onDragEnd, onDrop, dragging } = useFileDrop(props);

  return (
    <Box
      id="file-drop-zone"
      width="full"
      height="full"
      onDragEnter={onDragStart}
      onDragLeave={onDragEnd}
      onDrop={onDrop}
      onDragOver={(e) => e.preventDefault()}
      backgroundColor={dragging ? "gray.50" : ""}
      outline={dragging ? "2px var(--chakra-colors-red-200) dashed" : ""}
      outlineOffset={dragging ? "5px" : ""}
    >
      {children}
    </Box>
  );
}
