import { useToast } from "@chakra-ui/react";
import { DragEvent, useState } from "react";
import { Transforms } from "slate";
import { useSlate } from "slate-react";

import { assetUpload } from "src/api/openapi/assets";
import { Asset } from "src/api/openapi/schemas";

import { isSupportedImage } from "./utils";

export type Props = {
  onComplete: (asset: Asset) => void;
};

export function useFileDrop(props: Props) {
  const [dragging, setDragging] = useState(false);
  const editor = useSlate();
  const toast = useToast();

  function onDragStart() {
    setDragging(true);
  }

  function onDragEnd(e: DragEvent<HTMLDivElement>) {
    e.preventDefault();
    setDragging(false);
  }

  async function upload(f: File) {
    // TODO: Upload progress indicator...
    const asset = await assetUpload(f);

    props.onComplete(asset);

    return asset;
  }

  async function process(f: File) {
    const { url } = await upload(f);

    if (!isSupportedImage(f.type)) {
      throw new Error("Unsupported image format");
    }

    Transforms.insertNodes(editor, {
      type: "image",
      caption: url,
      link: url,
      children: [{ text: "" }],
    });
  }

  async function handleEvent(e: DragEvent<HTMLDivElement>) {
    if (e.dataTransfer.items) {
      await Promise.all(
        [...e.dataTransfer.items].map(async (item) => {
          if (item.kind === "file") {
            const file = item.getAsFile();
            if (file == null) return;

            await process(file);
          }
        })
      );
    } else {
      await Promise.all([...e.dataTransfer.files].map(process));
    }
  }

  async function onDrop(e: DragEvent<HTMLDivElement>) {
    e.preventDefault();

    try {
      await handleEvent(e);
    } catch (e: unknown) {
      toast({
        title: "Image upload failed",
        description: `${e}`,
      });
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
