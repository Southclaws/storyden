"use client";

import {
  type ClipboardEvent,
  type DragEvent,
  useEffect,
  useRef,
  useState,
} from "react";
import { toast } from "sonner";

import { Asset } from "@/api/openapi-schema";
import { useImageUpload } from "@/components/content/useImageUpload";
import { deriveError } from "@/utils/error";
import { generateXid } from "@/utils/xid";

const supportedImageTypes = new Set([
  "image/gif",
  "image/jpeg",
  "image/png",
  "image/webp",
]);

export type RobotChatImageAttachment = {
  id: string;
  filename: string;
  previewURL?: string;
  asset?: Asset;
};

export function useRobotChatImageAttachments() {
  const { upload } = useImageUpload();
  const [attachments, setAttachments] = useState<RobotChatImageAttachment[]>(
    [],
  );
  const [isDragging, setIsDragging] = useState(false);
  const dragCounter = useRef(0);
  const previewURLs = useRef(new Set<string>());

  useEffect(() => {
    const currentPreviewURLs = previewURLs.current;
    return () => {
      for (const previewURL of currentPreviewURLs) {
        URL.revokeObjectURL(previewURL);
      }
      currentPreviewURLs.clear();
    };
  }, []);

  function releasePreview(previewURL?: string) {
    if (!previewURL || !previewURLs.current.delete(previewURL)) {
      return;
    }
    URL.revokeObjectURL(previewURL);
  }

  async function uploadAttachment(
    attachment: RobotChatImageAttachment,
    file: File,
  ) {
    try {
      const asset = await upload(file, { filename: file.name });
      releasePreview(attachment.previewURL);
      setAttachments((current) =>
        current.map((item) =>
          item.id === attachment.id
            ? { ...item, asset, previewURL: undefined }
            : item,
        ),
      );
    } catch (error) {
      releasePreview(attachment.previewURL);
      setAttachments((current) =>
        current.filter((item) => item.id !== attachment.id),
      );
      toast.error(`Could not upload ${file.name}. ${deriveError(error)}`);
    }
  }

  function addFiles(files: File[] | FileList) {
    const selected = Array.from(files);
    const images = selected.filter((file) =>
      supportedImageTypes.has(file.type.toLowerCase()),
    );

    if (images.length !== selected.length) {
      toast.error("Images must be GIF, JPEG, PNG, or WebP files.");
    }

    const additions = images.map((file) => {
      const previewURL = URL.createObjectURL(file);
      previewURLs.current.add(previewURL);
      return {
        attachment: {
          id: generateXid(),
          filename: file.name,
          previewURL,
        } satisfies RobotChatImageAttachment,
        file,
      };
    });

    setAttachments((current) => [
      ...current,
      ...additions.map(({ attachment }) => attachment),
    ]);

    for (const { attachment, file } of additions) {
      void uploadAttachment(attachment, file);
    }
  }

  function handlePaste(event: ClipboardEvent<HTMLTextAreaElement>) {
    const files = Array.from(event.clipboardData.items).flatMap((item) => {
      if (item.kind !== "file" || !item.type.startsWith("image/")) {
        return [];
      }

      const file = item.getAsFile();
      return file ? [file] : [];
    });
    if (files.length === 0) {
      return;
    }

    event.preventDefault();
    addFiles(files);
  }

  function handleDragEnter(event: DragEvent<HTMLFormElement>) {
    if (!event.dataTransfer.types.includes("Files")) {
      return;
    }

    event.preventDefault();
    dragCounter.current += 1;
    setIsDragging(true);
  }

  function handleDragLeave(event: DragEvent<HTMLFormElement>) {
    if (!event.dataTransfer.types.includes("Files")) {
      return;
    }

    event.preventDefault();
    dragCounter.current = Math.max(0, dragCounter.current - 1);
    if (dragCounter.current === 0) {
      setIsDragging(false);
    }
  }

  function handleDragOver(event: DragEvent<HTMLFormElement>) {
    if (!event.dataTransfer.types.includes("Files")) {
      return;
    }

    event.preventDefault();
    event.dataTransfer.dropEffect = "copy";
  }

  function handleDrop(event: DragEvent<HTMLFormElement>) {
    if (!event.dataTransfer.types.includes("Files")) {
      return;
    }

    event.preventDefault();
    dragCounter.current = 0;
    setIsDragging(false);
    addFiles(event.dataTransfer.files);
  }

  function removeAttachment(id: string) {
    setAttachments((current) => {
      const removed = current.find((item) => item.id === id);
      releasePreview(removed?.previewURL);
      return current.filter((item) => item.id !== id);
    });
  }

  function clearAttachments() {
    setAttachments((current) => {
      for (const attachment of current) {
        releasePreview(attachment.previewURL);
      }
      return [];
    });
  }

  function restoreAssets(assets: Asset[]) {
    setAttachments((current) => [
      ...assets.map((asset) => ({
        id: generateXid(),
        filename: asset.filename,
        asset,
      })),
      ...current,
    ]);
  }

  const assets = attachments.flatMap((attachment) =>
    attachment.asset ? [attachment.asset] : [],
  );

  return {
    attachments,
    assets,
    isDragging,
    isUploading: assets.length !== attachments.length,
    addFiles,
    handlePaste,
    handleDragEnter,
    handleDragLeave,
    handleDragOver,
    handleDrop,
    removeAttachment,
    clearAttachments,
    restoreAssets,
  };
}
