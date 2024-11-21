import Image, { ImageOptions } from "@tiptap/extension-image";
import {
  NodeViewContent,
  NodeViewProps,
  NodeViewWrapper,
  ReactNodeViewRenderer,
  mergeAttributes,
} from "@tiptap/react";
import { Plugin } from "prosemirror-state";

import { Asset } from "src/api/openapi-schema";

import { css } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";

// NOTE: This is the name of the component that will be used in the HTML.
// It cannot be changed.
const COMPONENT_NAME = "img";

type Options = {
  handleFiles: (file: File[]) => Promise<Asset[]>;
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
          handleDOMEvents: {
            drop(view, event) {
              const hasFiles =
                event.dataTransfer &&
                event.dataTransfer.files &&
                event.dataTransfer.files.length;

              if (!hasFiles) {
                return;
              }

              const images = Array.from(event.dataTransfer.files).filter(
                (file) => /image/i.test(file.type),
              );

              if (images.length === 0) {
                return;
              }

              event.preventDefault();
              handleFiles?.(images);
            },

            paste(view, event) {
              const hasFiles =
                event.clipboardData &&
                event.clipboardData.files &&
                event.clipboardData.files.length;

              if (!hasFiles) {
                return;
              }

              const images = Array.from(event.clipboardData.files).filter(
                (file) => /image/i.test(file.type),
              );

              if (images.length === 0) {
                return;
              }

              event.preventDefault();
              handleFiles?.(images);
            },
          },
        },
      }),
    ];
  },
});
