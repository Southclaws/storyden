import { DragEvent, useState } from "react";
import { Transforms } from "slate";
import { useSlate } from "slate-react";

import { assetGetUploadURL } from "src/api/openapi/assets";
import { AssetGetUploadURLOKResponse } from "src/api/openapi/schemas";

export function useFileDrop() {
  const [dragging, setDragging] = useState(false);
  const editor = useSlate();

  function onDragStart() {
    setDragging(true);
  }

  function onDragEnd(e: DragEvent<HTMLDivElement>) {
    e.preventDefault();
    setDragging(false);
  }

  async function upload(f: File) {
    const { url } = await assetGetUploadURL();

    // TODO: Upload progress indicator...
    const response = await fetch(url, {
      credentials: "include",
      method: "POST",
      headers: { "Content-Type": "application/octet-stream" },
      body: f,
    });

    const json = (await response.json()) as AssetGetUploadURLOKResponse;

    return json.url;
  }

  async function process(f: File) {
    const url = await upload(f);

    Transforms.insertNodes(editor, {
      type: "image",
      caption: url,
      link: url,
      children: [{ text: "" }],
    });
  }

  async function onDrop(e: DragEvent<HTMLDivElement>) {
    e.preventDefault();

    if (e.dataTransfer.items) {
      [...e.dataTransfer.items].forEach((item) => {
        if (item.kind === "file") {
          const file = item.getAsFile();
          if (file == null) return;

          process(file);
        }
      });
    } else {
      [...e.dataTransfer.files].forEach(process);
    }

    setDragging(false);
  }

  return {
    onDragStart,
    onDragEnd,
    onDrop,
    dragging,
  };
}
