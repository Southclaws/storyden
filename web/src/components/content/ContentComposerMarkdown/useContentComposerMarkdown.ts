import { useEffect, useRef, useState } from "react";

import { handle } from "@/api/client";
import { getAssetURL } from "@/utils/asset";
import { htmlToMarkdown, markdownToHTML } from "@/utils/markdown";

import { ContentComposerProps } from "../composer-props";
import {
  getImageFiles,
  hasImageFile,
  isSupportedImage,
  useImageUpload,
} from "../useImageUpload";

const PLACEHOLDER_SEPARATOR = "\n\n";
const ERROR_UNSUPPORTED_FILE_TYPE = "File type not supported";

export function useContentComposerMarkdown(props: ContentComposerProps) {
  const [value, setValue] = useState(() => {
    if (props.initialValue) {
      return htmlToMarkdown(props.initialValue);
    }
    return "";
  });
  const { upload } = useImageUpload();
  const [previewHTML, setPreviewHTML] = useState<string>("");
  const [showPreview, setShowPreview] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const [isDragError, setIsDragError] = useState(false);
  const [dragErrorMessage, setDragErrorMessage] = useState("");
  const [dragFileCount, setDragFileCount] = useState(0);
  const [uploadingCount, setUploadingCount] = useState(0);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const uploadCounterRef = useRef(0);
  const dragCounterRef = useRef(0);

  useEffect(() => {
    if (props.resetKey) {
      setValue("");
      setPreviewHTML("");
    }
  }, [props.resetKey]);

  useEffect(() => {
    if (props.disabled || showPreview) return;

    const textarea = textareaRef.current;
    if (!textarea) return;

    textarea.style.height = "auto";
    textarea.style.height = `${textarea.scrollHeight}px`;
  }, [
    value,
    // When the editor is disabled, we don't need to adjust the height.
    props.disabled,
    // We depend on showPreview here because when preview is disabled, we want
    // to recalculate the height of the textarea.
    showPreview,
  ]);

  async function handleBufferChange(e: React.ChangeEvent<HTMLTextAreaElement>) {
    const markdownRaw = e.target.value;

    setValue(markdownRaw);

    const html = await markdownToHTML(markdownRaw);

    const isEmpty = markdownRaw.trim().length === 0 || html.trim().length === 0;

    if (props.onChange) {
      props.onChange(html, isEmpty);
    }
  }

  async function handleTogglePreview() {
    if (!showPreview) {
      const html = await markdownToHTML(value);
      setPreviewHTML(html);
    }
    setShowPreview(!showPreview);
  }

  function getDragOverlayMessage() {
    if (isDragError) {
      return dragErrorMessage;
    }
    return dragFileCount === 1
      ? "Drop 1 file to upload"
      : `Drop ${dragFileCount} files to upload`;
  }

  async function handlePaste(e: React.ClipboardEvent<HTMLTextAreaElement>) {
    const items = Array.from(e.clipboardData.items);
    const imageFiles: File[] = [];

    for (const item of items) {
      if (isSupportedImage(item.type)) {
        const file = item.getAsFile();
        if (file) {
          imageFiles.push(file);
        }
      }
    }

    if (imageFiles.length > 0) {
      e.preventDefault();
      await handleMultipleImageUploads(imageFiles);
    }
  }

  async function handleDrop(e: React.DragEvent<HTMLTextAreaElement>) {
    e.preventDefault();

    dragCounterRef.current = 0;
    setIsDragging(false);
    setIsDragError(false);
    setDragErrorMessage("");
    setDragFileCount(0);

    const imageFiles = getImageFiles(e.dataTransfer.files);

    if (imageFiles.length > 0) {
      await handleMultipleImageUploads(imageFiles);
    }
  }

  async function handleMultipleImageUploads(files: File[]) {
    const textarea = textareaRef.current;
    if (!textarea) return;

    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;

    const placeholders: Array<{ file: File; placeholder: string; id: number }> =
      [];
    let combinedPlaceholders = "";

    for (const file of files) {
      uploadCounterRef.current += 1;
      const uploadId = uploadCounterRef.current;
      const placeholder = `![Uploading image ${uploadId}...]()`;
      placeholders.push({ file, placeholder, id: uploadId });
      combinedPlaceholders += placeholder + PLACEHOLDER_SEPARATOR;
    }

    const newValue =
      value.substring(0, start) + combinedPlaceholders + value.substring(end);
    setValue(newValue);

    setTimeout(() => {
      textarea.focus();
      const newCursorPos = start + combinedPlaceholders.length;
      textarea.setSelectionRange(newCursorPos, newCursorPos);
    }, 0);

    await Promise.all(
      placeholders.map(({ file, placeholder }) =>
        handleImageUploadWithPlaceholder(file, placeholder, start),
      ),
    );
  }

  async function handleImageUploadWithPlaceholder(
    file: File,
    placeholder: string,
    insertPosition: number,
  ) {
    setUploadingCount((prev) => prev + 1);

    await handle(
      async () => {
        const asset = await upload(file, { filename: file.name });

        const imageUrl = getAssetURL(asset.path);
        const markdownImage = `![${file.name}](${imageUrl})`;

        setValue((currentValue) => {
          const placeholderIndex = currentValue.indexOf(
            placeholder,
            insertPosition,
          );
          if (placeholderIndex === -1) {
            return currentValue;
          }

          const newValue =
            currentValue.substring(0, placeholderIndex) +
            markdownImage +
            currentValue.substring(placeholderIndex + placeholder.length);

          markdownToHTML(newValue).then((html) => {
            const isEmpty =
              newValue.trim().length === 0 || html.trim().length === 0;
            if (props.onChange) {
              props.onChange(html, isEmpty);
            }
          });

          return newValue;
        });

        props.onAssetUpload?.(asset);
      },
      {
        cleanup: async () => {
          setUploadingCount((prev) => prev - 1);
        },
      },
    );
  }

  function handleDragOver(e: React.DragEvent<HTMLTextAreaElement>) {
    e.preventDefault();
  }

  function handleDragEnter(e: React.DragEvent<HTMLTextAreaElement>) {
    const items = Array.from(e.dataTransfer.items);
    const hasFile = items.some((item) => item.kind === "file");

    if (!hasFile) {
      return;
    }

    e.preventDefault();
    dragCounterRef.current += 1;
    setIsDragging(true);

    const hasImage = hasImageFile(e.dataTransfer.items);
    const imageCount = items.filter((item) =>
      isSupportedImage(item.type),
    ).length;

    if (!hasImage) {
      setIsDragError(true);
      setDragErrorMessage(ERROR_UNSUPPORTED_FILE_TYPE);
    } else {
      setIsDragError(false);
      setDragErrorMessage("");
    }

    setDragFileCount(imageCount);
  }

  function handleDragLeave() {
    dragCounterRef.current -= 1;
    if (dragCounterRef.current === 0) {
      setIsDragging(false);
      setIsDragError(false);
      setDragErrorMessage("");
      setDragFileCount(0);
    }
  }

  return {
    value,
    previewHTML,
    showPreview,
    isDragging,
    isDragError,
    uploadingCount,
    textareaRef,
    getDragOverlayMessage,

    // Handlers
    handleBufferChange,
    handleTogglePreview,
    handlePaste,
    handleDrop,
    handleDragOver,
    handleDragEnter,
    handleDragLeave,
  };
}
