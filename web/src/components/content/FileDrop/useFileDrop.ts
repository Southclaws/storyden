"use client";

import { DragEvent, useState } from "react";

import { assetUpload } from "src/api/openapi/assets";
import { Asset } from "src/api/openapi/schemas";

import { isSupportedImage } from "./utils";

export type Props = {
  disabled?: boolean;
  onComplete?: (asset: Asset) => void;
};

export function useFileUpload() {
  async function upload(f: File) {
    if (!isSupportedImage(f.type)) {
      throw new Error("Unsupported image format");
    }

    const asset = await assetUpload(f);

    return asset;
  }

  return {
    upload,
  };
}

export function useFileDrop(props: Props) {
  const [dragging, setDragging] = useState(false);
  const { upload } = useFileUpload();

  async function handleUpload(f: File) {
    const asset = await upload(f);
    props.onComplete?.(asset);
  }

  function onDragStart() {
    if (props.disabled) return;

    setDragging(true);
  }

  function onDragEnd(e: DragEvent<HTMLDivElement>) {
    if (props.disabled) return;

    e.preventDefault();
    setDragging(false);
  }

  async function handleEvent(e: DragEvent<HTMLDivElement>) {
    if (props.disabled) return;

    if (e.dataTransfer.items) {
      await Promise.all(
        [...e.dataTransfer.items].map(async (item) => {
          if (item.kind === "file") {
            const file = item.getAsFile();
            if (file == null) return;

            await handleUpload(file);
          }
        }),
      );
    } else {
      await Promise.all([...e.dataTransfer.files].map(handleUpload));
    }
  }

  async function onDrop(e: DragEvent<HTMLDivElement>) {
    if (props.disabled) return;

    e.preventDefault();

    try {
      await handleEvent(e);
    } catch (e: unknown) {
      console.error(e);
    } finally {
      setDragging(false);
    }
  }

  return {
    onDragStart,
    onDragEnd,
    onDrop,
    dragging,
  };
}
