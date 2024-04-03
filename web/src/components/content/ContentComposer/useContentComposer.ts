"use client";

import Mention from "@tiptap/extension-mention";
import Placeholder from "@tiptap/extension-placeholder";
import { useEditor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";

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

  const editor = useEditor({
    editorProps: {
      attributes: {
        class: css({
          height: "full",
          width: "full",
          boxShadow: "md",
          padding: "2",
          borderRadius: "sm",
        }),
      },
    },
    extensions: [
      StarterKit,
      ImageExtended.configure({
        allowBase64: false,
        HTMLAttributes: {
          class: css({ borderRadius: "md" }),
        },
        handleFileUpload: upload,
      }),
      Placeholder.configure({
        placeholder: "Write your heart out...",
        includeChildren: true,
        showOnlyCurrent: false,
        considerAnyAsEmpty: true,
      }),
      Mention.configure({
        HTMLAttributes: {
          class: "mention",
        },
        suggestion: {
          render: () => ({
            onStart: (props) => {
              console.log("start", props.clientRect?.());
            },
          }),
        },
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

  return {
    editor,
    handlers: {
      handleBold,
      handleItalic,
      handleStrike,
    },
  };
}
