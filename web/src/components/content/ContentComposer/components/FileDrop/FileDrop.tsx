"use client";

import { PropsWithChildren } from "react";

import { Box } from "@/styled-system/jsx";

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
      style={{
        backgroundColor: dragging ? "gray.50" : undefined,
        outline: dragging
          ? "2px var(--chakra-colors-red-200) dashed"
          : undefined,
        outlineOffset: dragging ? "0.5" : undefined,
      }}
    >
      {children}
    </Box>
  );
}
