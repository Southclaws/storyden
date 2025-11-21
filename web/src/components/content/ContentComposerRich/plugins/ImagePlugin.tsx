import Image, { ImageOptions } from "@tiptap/extension-image";
import {
  NodeViewContent,
  NodeViewProps,
  NodeViewWrapper,
  ReactNodeViewRenderer,
  mergeAttributes,
} from "@tiptap/react";
import { EditorView } from "prosemirror-view";
import { Plugin } from "prosemirror-state";

import { Asset } from "src/api/openapi-schema";

import { css } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";

// NOTE: This is the name of the component that will be used in the HTML.
// It cannot be changed.
const COMPONENT_NAME = "img";

type Options = {
  handleFiles: (view: EditorView, files: File[]) => Promise<Asset[]>;
};

function Component(props: NodeViewProps) {
  return (
    <NodeViewWrapper
      className={css({
        cursor: "pointer",
      })}
    >
      <styled.img borderRadius="md" {...props.node.attrs} />
    </NodeViewWrapper>
  );
}

export const ImageExtended = Image.extend<ImageOptions & Options>({
  content: "inline*",
  addOptions() {
    return {
      ...this.parent?.(),
    };
  },
  addNodeView() {
    return ReactNodeViewRenderer(Component);
  },
  parseHTML() {
    return [
      {
        tag: COMPONENT_NAME,
      },
    ];
  },
  renderHTML({ HTMLAttributes }) {
    return [COMPONENT_NAME, mergeAttributes(HTMLAttributes), 0];
  },
  addProseMirrorPlugins() {
    const handleFiles = this.options.handleFiles;
    return [
      new Plugin({
        props: {
          handlePaste(view, event) {
            if (!event.clipboardData) {
              return false;
            }

            const files: File[] = [];

            // Use "items"
            if (event.clipboardData.items?.length) {
              for (const item of event.clipboardData.items) {
                if (item.kind === "file") {
                  const file = item.getAsFile();
                  if (file) {
                    files.push(file);
                  }
                }
              }
            }

            const images = files.filter((file) => /image/i.test(file.type));

            if (images.length === 0) {
              return false;
            }

            event.preventDefault();
            handleFiles(view, images);
            return true;
          },
          handleDOMEvents: {
            drop(view, event) {
              const hasFiles =
                event.dataTransfer &&
                event.dataTransfer.files &&
                event.dataTransfer.files.length;

              if (!hasFiles) {
                return false;
              }

              const images = Array.from(event.dataTransfer.files).filter(
                (file) => /image/i.test(file.type),
              );

              if (images.length === 0) {
                return false;
              }

              event.preventDefault();
              handleFiles(view, images);
              return true;
            },
          },
        },
      }),
    ];
  },
});
