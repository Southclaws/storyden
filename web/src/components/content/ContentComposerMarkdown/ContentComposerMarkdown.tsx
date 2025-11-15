import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useRef, useState } from "react";
import Markdown from "react-markdown";

import { handle } from "@/api/client";
import { assetUpload } from "@/api/openapi-client/assets";
import { Spinner } from "@/components/ui/Spinner";
import { IconButton } from "@/components/ui/icon-button";
import { EditIcon } from "@/components/ui/icons/Edit";
import { ShowIcon } from "@/components/ui/icons/ShowIcon";
import { Switch } from "@/components/ui/switch";
import { Box, HStack, LStack, styled } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";
import { htmlToMarkdown, markdownToHTML } from "@/utils/markdown";

import { ContentComposerProps } from "../composer-props";
import {
  getImageFiles,
  hasImageFile,
  isSupportedImage,
} from "../useImageUpload";

const PLACEHOLDER_SEPARATOR = "\n\n";
const ERROR_UNSUPPORTED_FILE_TYPE = "File type not supported";

export function ContentComposerMarkdown(props: ContentComposerProps) {
  const [value, setValue] = useState(() => {
    if (props.initialValue) {
      return htmlToMarkdown(props.initialValue);
    }
    return "";
  });
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
    }
  }, [props.resetKey]);

  useEffect(() => {
    if (props.disabled) return;

    const textarea = textareaRef.current;
    if (!textarea) return;

    const parent = textarea.parentElement;
    if (!parent) return;

    const resizeObserver = new ResizeObserver(() => {
      textarea.style.height = "0px";
      textarea.style.height = `${textarea.scrollHeight}px`;
    });

    resizeObserver.observe(parent);

    textarea.style.height = "0px";
    textarea.style.height = `${textarea.scrollHeight}px`;

    return () => {
      resizeObserver.disconnect();
    };
  }, [value, props.disabled]);

  async function onChange(e: React.ChangeEvent<HTMLTextAreaElement>) {
    const markdownRaw = e.target.value;

    setValue(markdownRaw);

    const html = await markdownToHTML(markdownRaw);

    const isEmpty = markdownRaw.trim().length === 0 || html.trim().length === 0;

    if (props.onChange) {
      props.onChange(html, isEmpty);
    }
  }

  function handleTogglePreview() {
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
        const asset = await assetUpload(file, { filename: file.name });
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

  if (props.disabled) {
    return (
      <LStack position="relative" minHeight="8" flex="1">
        <Markdown className="typography">{value}</Markdown>
      </LStack>
    );
  }

  return (
    <LStack position="relative" minHeight="8" flex="1">
      <EditorTools
        showPreview={showPreview}
        onChange={handleTogglePreview}
        workingCount={uploadingCount}
      />

      {showPreview ? (
        <>
          {value ? (
            <Markdown className="typography">{value}</Markdown>
          ) : (
            <styled.p height="14" color="fg.muted" fontStyle="italic">
              empty...
            </styled.p>
          )}
        </>
      ) : (
        <>
          <styled.textarea
            ref={textareaRef}
            onChange={onChange}
            onPaste={handlePaste}
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onDragEnter={handleDragEnter}
            onDragLeave={handleDragLeave}
            value={value}
            lineHeight="relaxed"
            w="full"
            minHeight="0"
            resize="none"
            appearance="none"
            border="none"
            outline="none"
            color="fg.default"
            fontSize="md"
            transitionDuration="normal"
            transitionTimingFunction="default"
            _placeholder={{
              color: "fg.default",
            }}
            style={{
              border: "none",
              transitionProperty: "border-color, border-width",
              overflow: "hidden",
            }}
            placeholder="Write your heart out..."
          />
          {isDragging && (
            <Box
              position="absolute"
              top="0"
              left="0"
              right="0"
              bottom="0"
              pointerEvents="none"
              display="flex"
              alignItems="center"
              justifyContent="center"
              backgroundColor="bg.emphasized"
              borderWidth="medium"
              borderStyle="dashed"
              borderColor={isDragError ? "border.error" : "accent.default"}
              borderRadius="md"
              style={{ opacity: 0.95 }}
              role="status"
              aria-live="polite"
              aria-label={getDragOverlayMessage()}
            >
              <styled.div
                fontSize="sm"
                fontWeight="medium"
                color={isDragError ? "fg.error" : "accent.default"}
                display="flex"
                flexDirection="column"
                alignItems="center"
                gap="2"
              >
                <span>{getDragOverlayMessage()}</span>
              </styled.div>
            </Box>
          )}
        </>
      )}
    </LStack>
  );
}

function EditorTools({
  showPreview,
  onChange,
  workingCount,
}: {
  showPreview: boolean;
  onChange: () => void;
  workingCount: number;
}) {
  const [isHovered, setIsHovered] = useState(false);

  const isWorking = workingCount > 0;

  // Logic: reveal the animated presence container if working or hovered, this
  // allows the spinner to appear on its own when the editor is working. When
  // hovered, we also include the preview switch. If the user hovers while the
  // editor is working, the preview switch appears alongside the spinner.
  const reveal = isWorking || isHovered;

  return (
    <Box
      position="absolute"
      right="0"
      p="1"
      opacity={isHovered ? "full" : "5"}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      cursor="pointer"
      backgroundColor="bg.subtle"
      backdropBlur="frosted"
      backdropFilter="auto"
      borderRadius="md"
      transition="all"
    >
      <HStack gap="2">
        <AnimatePresence>
          {reveal && (
            <motion.div
              initial={{ width: 0, opacity: 0 }}
              animate={{ width: "auto", opacity: 1 }}
              exit={{ width: 0, opacity: 0 }}
              transition={{ duration: 0.2, ease: "easeInOut" }}
              style={{ overflow: "hidden" }}
            >
              <HStack gap="2">
                {isWorking && (
                  <HStack gap="1">
                    <Spinner w="4" h="4" />
                    {workingCount > 1 && (
                      <styled.span fontSize="xs" color="fg.muted">
                        {workingCount}
                      </styled.span>
                    )}
                  </HStack>
                )}
                {isHovered && (
                  <Switch size="sm" checked={showPreview} onClick={onChange}>
                    Preview
                  </Switch>
                )}
              </HStack>
            </motion.div>
          )}
        </AnimatePresence>

        <IconButton type="button" variant="ghost" size="xs" onClick={onChange}>
          {showPreview ? <EditIcon w="4" /> : <ShowIcon w="4" />}
        </IconButton>
      </HStack>
    </Box>
  );
}
