"use client";

import { FocusClasses } from "@tiptap/extension-focus";
import { Link } from "@tiptap/extension-link";
import Placeholder from "@tiptap/extension-placeholder";
import { generateHTML, generateJSON } from "@tiptap/html";
import { EditorView } from "@tiptap/pm/view";
import { useEditor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import { ChangeEvent, useEffect, useId, useRef, useState } from "react";

import { Asset } from "src/api/openapi-schema";

import { handle } from "@/api/client";
import { css } from "@/styled-system/css";
import { getAssetURL } from "@/utils/asset";

import { ContentComposerProps } from "../composer-props";
import {
  hasImageFile,
  isSupportedImage,
  useImageUpload,
} from "../useImageUpload";

import { ImageExtended } from "./plugins/ImagePlugin";

const ERROR_UNSUPPORTED_FILE_TYPE = "File type not supported";

export type Block = "p" | "h1" | "h2" | "h3" | "h4" | "h5" | "h6";

export function useContentComposer(props: ContentComposerProps) {
  const { upload } = useImageUpload();
  const [uploadingCount, setUploadingCount] = useState(0);
  const [isDragging, setIsDragging] = useState(false);
  const [isDragError, setIsDragError] = useState(false);
  const [dragErrorMessage, setDragErrorMessage] = useState("");
  const [dragFileCount, setDragFileCount] = useState(0);
  const dragCounterRef = useRef(0);

  const extensions = [
    StarterKit,
    FocusClasses,
    Link.configure({
      // Disable navigation when clicking links in the editor if it's active.
      openOnClick: props.disabled ? true : false,
    }).extend({
      inclusive: false,
    }),
    ImageExtended.configure({
      allowBase64: false,
      HTMLAttributes: {
        class: css({ borderRadius: "md" }),
      },
      handleFiles,
    }),
    Placeholder.configure({
      placeholder: props.placeholder ?? "Write your heart out...",
      includeChildren: true,
      showOnlyCurrent: false,
    }),
  ];

  // This is for the initial server render.
  const initialValueJSON = generateJSON(
    props.initialValue ?? "<p></p>",
    extensions,
  );
  const initialValueHTML = generateHTML(initialValueJSON, extensions);

  // Each editor needs a unique ID for the menu's file upload input ID.
  const uniqueID = useId();

  const editor = useEditor({
    immediatelyRender: false,
    editorProps: {
      attributes: {
        "data-editor-id": uniqueID,
        class: css({
          height: "full",
          width: "full",
        }),
      },
    },
    extensions,
    content: props.initialValue ?? "<p></p>",
    onUpdate: ({ editor }) => {
      const html = editor.getHTML();
      props.onChange?.(html, editor.isEmpty);
    },
  });

  // This is a huge hack but it means the composer doesn't need to be made into
  // a controlled component. Baiscally, if the resetKey changes, we reset the
  // content of the editor to the initial value or empty paragraph. Hacky? Yes.
  useEffect(() => {
    if (!editor) {
      return;
    }

    if (!props.resetKey) {
      if (props.value) {
        editor.commands.setContent(props.value);
      }
      return;
    }

    editor.commands.setContent(props.initialValue ?? "<p></p>");
  }, [editor, props.initialValue, props.value, props.resetKey]);

  useEffect(() => {
    if (!editor) {
      return;
    }

    editor.setEditable(!props.disabled, false);
  }, [editor, props.disabled]);

  async function handleFiles(view: EditorView, files: File[]) {
    if (!view) {
      return [];
    }

    const { state } = view;
    const { selection } = state;
    const { schema } = view.state;
    const imageNode = schema.nodes?.["image"];

    if (!imageNode) {
      return [];
    }

    const assets: Asset[] = [];
    for (const f of files) {
      setUploadingCount((prev) => prev + 1);

      await handle(
        async () => {
          const asset = await upload(f);

          const node = imageNode.create({ src: getAssetURL(asset.path) });
          const transaction = view.state.tr.insert(selection.$head.pos, node);
          view.dispatch(transaction);

          assets.push(asset);
          props.onAssetUpload?.(asset);
        },
        {
          cleanup: async () => {
            setUploadingCount((prev) => prev - 1);
          },
        },
      );
    }

    return assets;
  }

  async function handleFileUpload(e: ChangeEvent<HTMLInputElement>) {
    if (!e.currentTarget.files || !editor) {
      return;
    }

    const images = Array.from(e.currentTarget.files).filter((file) =>
      /image/i.test(file.type),
    );

    await handleFiles(editor.view, images);
  }

  function handleBlockType(kind: Block) {
    switch (kind) {
      case "p":
        editor?.chain().focus().setParagraph().run();
        break;
      case "h1":
        editor?.chain().focus().setHeading({ level: 1 }).run();
        break;
      case "h2":
        editor?.chain().focus().setHeading({ level: 2 }).run();
        break;

      case "h3":
        editor?.chain().focus().setHeading({ level: 3 }).run();
        break;

      case "h4":
        editor?.chain().focus().setHeading({ level: 4 }).run();
        break;

      case "h5":
        editor?.chain().focus().setHeading({ level: 5 }).run();
        break;

      case "h6":
        editor?.chain().focus().setHeading({ level: 6 }).run();
        break;
    }
  }

  function getBlockType(): Block | null {
    if (editor?.isActive("paragraph")) return "p";
    if (editor?.isActive("heading", { level: 1 })) return "h1";
    if (editor?.isActive("heading", { level: 2 })) return "h2";
    if (editor?.isActive("heading", { level: 3 })) return "h3";
    if (editor?.isActive("heading", { level: 4 })) return "h4";
    if (editor?.isActive("heading", { level: 5 })) return "h5";
    if (editor?.isActive("heading", { level: 6 })) return "h6";

    return "p";
  }

  function getDragOverlayMessage() {
    if (isDragError) {
      return dragErrorMessage;
    }
    return dragFileCount === 1
      ? "Drop 1 file to upload"
      : `Drop ${dragFileCount} files to upload`;
  }

  function handleDragOver(e: React.DragEvent) {
    e.preventDefault();
  }

  function handleDragEnter(e: React.DragEvent) {
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

  async function handleDrop(e: React.DragEvent) {
    e.preventDefault();

    dragCounterRef.current = 0;
    setIsDragging(false);
    setIsDragError(false);
    setDragErrorMessage("");
    setDragFileCount(0);

    if (!editor) {
      return;
    }

    const files = Array.from(e.dataTransfer.files);
    const images = files.filter((file) => /image/i.test(file.type));

    if (images.length > 0) {
      await handleFiles(editor.view, images);
    }
  }

  return {
    editor,
    uniqueID,
    initialValueHTML,
    uploadingCount,
    isDragging,
    isDragError,
    getDragOverlayMessage,
    handlers: {
      handleFileUpload,
      handleDragOver,
      handleDragEnter,
      handleDragLeave,
      handleDrop,
    },
    format: {
      text: {
        active: getBlockType(),
        set: handleBlockType,
      },

      // Marks
      bold: {
        isActive: editor?.isActive("bold") ?? false,
        isDisabled: editor?.can().toggleBold() === false,
        toggle: () => editor?.chain().focus().toggleBold().run(),
      },
      italic: {
        isActive: editor?.isActive("italic") ?? false,
        isDisabled: editor?.can().toggleItalic() === false,
        toggle: () => editor?.chain().focus().toggleItalic().run(),
      },
      strike: {
        isActive: editor?.isActive("strike") ?? false,
        isDisabled: editor?.can().toggleStrike() === false,
        toggle: () => editor?.chain().focus().toggleStrike().run(),
      },
      code: {
        isActive: editor?.isActive("code") ?? false,
        isDisabled: editor?.can().toggleCode() === false,
        toggle: () => editor?.chain().focus().toggleCode().run(),
      },

      // Blocks
      blockquote: {
        isActive: editor?.isActive("blockquote") ?? false,
        isDisabled: editor?.can().toggleBlockquote() === false,
        toggle: () => editor?.chain().focus().toggleBlockquote().run(),
      },
      pre: {
        isActive: editor?.isActive("codeBlock") ?? false,
        isDisabled: editor?.can().toggleCodeBlock() === false,
        toggle: () => editor?.chain().focus().toggleCodeBlock().run(),
      },
      bulletList: {
        isActive: editor?.isActive("bulletList") ?? false,
        isDisabled: editor?.can().toggleBulletList() === false,
        toggle: () => editor?.chain().focus().toggleBulletList().run(),
      },
      orderedList: {
        isActive: editor?.isActive("orderedList") ?? false,
        isDisabled: editor?.can().toggleOrderedList() === false,
        toggle: () => editor?.chain().focus().toggleOrderedList().run(),
      },
    },
  };
}
