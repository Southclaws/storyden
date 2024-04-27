"use client";

import { FocusClasses } from "@tiptap/extension-focus";
import Placeholder from "@tiptap/extension-placeholder";
import { useEditor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import { ChangeEvent } from "react";

import { useImageUpload } from "../useImageUpload";

import { css } from "@/styled-system/css";

import { ImageExtended } from "./plugins/ImagePlugin";

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

    const assets = [];
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

  function handleBold() {
    editor?.chain().focus().toggleBold().run();
  }

  function handleItalic() {
    editor?.chain().focus().toggleItalic().run();
  }

  function handleStrike() {
    editor?.chain().focus().toggleStrike().run();
  }

  async function handleFileUpload(e: ChangeEvent<HTMLInputElement>) {
    if (!e.currentTarget.files) {
      return;
    }

    const images = Array.from(e.currentTarget.files).filter((file) =>
      /image/i.test(file.type),
    );

    await handleFiles(images);
  }

  return {
    editor,
    handlers: {
      handleBold,
      handleItalic,
      handleStrike,
      handleFileUpload,
    },
  };
}
