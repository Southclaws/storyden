"use client";

import { Mark } from "@tiptap/core";
import { FocusClasses } from "@tiptap/extension-focus";
import { Link } from "@tiptap/extension-link";
import { generateJSON } from "@tiptap/html";
import { EditorContent, useEditor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import { useId, useMemo } from "react";

import { countDiffMarks, diffTipTapJSON } from "@/lib/content/diff";
import { css } from "@/styled-system/css";
import { cx } from "@/styled-system/css/index.mjs";
import { LStack } from "@/styled-system/jsx";

import { ImageExtended } from "../ContentComposerRich/plugins/ImagePlugin";
import { LinkPreview } from "../ContentComposerRich/plugins/LinkPreviewPlugin";

import "./diff.css";

export type DiffViewerInlineProps = {
  className?: string;
  originalHTML: string;
  modifiedHTML: string;
  showLegend?: boolean;
};

// Custom TipTap marks for insert/delete elements.

const DiffInsertion = Mark.create({
  name: "diffInsertion",

  parseHTML() {
    return [{ tag: "ins" }, { tag: "span[data-diff='insertion']" }];
  },

  renderHTML() {
    return [
      "span",
      {
        class: "diff-insertion",
        "data-diff": "insertion",
      },
      0,
    ];
  },

  addAttributes() {
    return {
      "data-diff": {
        default: "insertion",
        parseHTML: (element) => element.getAttribute("data-diff"),
        renderHTML: (attributes) => {
          return { "data-diff": attributes["data-diff"] };
        },
      },
    };
  },
});

const DiffDeletion = Mark.create({
  name: "diffDeletion",

  parseHTML() {
    return [{ tag: "del" }, { tag: "span[data-diff='deletion']" }];
  },

  renderHTML() {
    return [
      "span",
      {
        class: "diff-deletion",
        "data-diff": "deletion",
      },
      0,
    ];
  },

  addAttributes() {
    return {
      "data-diff": {
        default: "deletion",
        parseHTML: (element) => element.getAttribute("data-diff"),
        renderHTML: (attributes) => {
          return { "data-diff": attributes["data-diff"] };
        },
      },
    };
  },
});

/**
 * DiffViewerInline component displays an inline diff view with highlighting.
 */
export function ContentDiffView({
  className,
  originalHTML,
  modifiedHTML,
}: DiffViewerInlineProps) {
  const uniqueID = useId();

  // NOTE: The extensions here MUST match the extensions in the composer editor.
  const extensions = useMemo(
    () => [
      StarterKit,
      FocusClasses,
      DiffInsertion,
      DiffDeletion,
      LinkPreview,
      Link.configure({
        openOnClick: false,
      }).extend({
        inclusive: false,
        parseHTML() {
          return [
            {
              tag: "a[href]:not([data-display])",
              getAttrs: (dom) => {
                const href = dom.getAttribute("href");
                return href ? { href } : false;
              },
            },
          ];
        },
      }),
      ImageExtended.configure({
        allowBase64: false,
        HTMLAttributes: {
          class: css({ borderRadius: "md" }),
        },
        handleFiles: async () => [],
        handleRetry: () => {},
        handleCancel: () => {},
      }),
    ],
    [],
  );

  const { mergedDoc, stats } = useMemo(() => {
    const originalDoc = generateJSON(originalHTML, extensions);
    const modifiedDoc = generateJSON(modifiedHTML, extensions);

    const mergedDoc = diffTipTapJSON(originalDoc, modifiedDoc);

    const stats = countDiffMarks(mergedDoc);

    return { mergedDoc, stats };
  }, [originalHTML, modifiedHTML, extensions]);

  // TODO: Expose stats as an imperative handle.

  const editor = useEditor({
    immediatelyRender: false,
    editable: false,
    extensions,
    content: mergedDoc,
  });

  return (
    <LStack
      id={`rich-text-editor-${uniqueID}`}
      containerType="inline-size"
      className={cx("typography", className)}
      position="relative"
      w="full"
      gap="1"
      minHeight="8"
      // onDragOver={handlers.handleDragOver}
      // onDragEnter={handlers.handleDragEnter}
      // onDragLeave={handlers.handleDragLeave}
      // onDrop={handlers.handleDrop}
    >
      <div
        id={`editor-content-${uniqueID}`}
        className={css({
          height: "full",
          width: "full",
        })}
        suppressHydrationWarning
      >
        <EditorContent editor={editor} />
      </div>
    </LStack>
  );
}
