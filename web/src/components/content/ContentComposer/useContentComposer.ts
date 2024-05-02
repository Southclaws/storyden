"use client";

import { FocusClasses } from "@tiptap/extension-focus";
import Placeholder from "@tiptap/extension-placeholder";
import { useEditor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import { ChangeEvent } from "react";

import { Asset } from "src/api/openapi/schemas";

import { useImageUpload } from "../useImageUpload";

import { css } from "@/styled-system/css";

import { ImageExtended } from "./plugins/ImagePlugin";

export type Block = "p" | "h1" | "h2" | "h3" | "h4" | "h5" | "h6";

export type Props = {
  disabled?: boolean;
  initialValue?: string;
  onChange?: (value: string) => void;
};

export function useContentComposer(props: Props) {
  const { upload } = useImageUpload();

  async function handleFiles(files: File[]) {
    if (!editor) {
      return [];
    }

    const { view } = editor;
    const { state } = view;
    const { selection } = state;
    const { schema } = view.state;
    const imageNode = schema.nodes?.["image"];

    if (!imageNode) {
      return [];
    }

    const assets: Asset[] = [];
    for (const f of files) {
      const asset = await upload(f);

      const node = imageNode.create({ src: asset.url });
      const transaction = view.state.tr.insert(selection.$head.pos, node);
      view.dispatch(transaction);

      assets.push(asset);
    }

    return assets;
  }

  const editor = useEditor({
    editorProps: {
      attributes: {
        class: css({
          height: "full",
          width: "full",
        }),
      },
    },
    extensions: [
      StarterKit,
      FocusClasses,
      ImageExtended.configure({
        allowBase64: false,
        HTMLAttributes: {
          class: css({ borderRadius: "md" }),
        },
        handleFiles,
      }),
      Placeholder.configure({
        placeholder: "Write your heart out...",
        includeChildren: true,
        showOnlyCurrent: false,
        considerAnyAsEmpty: true,
      }),
    ],
    content: props.initialValue ?? "<p></p>",
    onUpdate: ({ editor }) => {
      const html = editor.getHTML();
      props.onChange?.(html);
    },
  });

  async function handleFileUpload(e: ChangeEvent<HTMLInputElement>) {
    if (!e.currentTarget.files) {
      return;
    }

    const images = Array.from(e.currentTarget.files).filter((file) =>
      /image/i.test(file.type),
    );

    await handleFiles(images);
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

    return null;
  }

  return {
    editor,
    handlers: {
      handleFileUpload,
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
