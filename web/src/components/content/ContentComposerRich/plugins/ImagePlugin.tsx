import type { ImageOptions } from "@tiptap/extension-image";
import { Plugin, PluginKey } from "@tiptap/pm/state";
import type { EditorView } from "@tiptap/pm/view";
import {
  type NodeViewProps,
  NodeViewWrapper,
  ReactNodeViewRenderer,
} from "@tiptap/react";

import type { Asset } from "@/api/openapi-schema";

import { ContentImage } from "../../ContentRenderer/ContentImage";
import { ImageNode } from "../../ContentRenderer/ImageNode";

// Track upload positions without repeatedly scanning the editor document.
export const uploadPositionsKey = new PluginKey<Map<string, number>>(
  "uploadPositions",
);

type Options = {
  handleFiles?: (view: EditorView, files: File[]) => Promise<Asset[]>;
  handleRetry?: (view: EditorView, uploadId: string) => void;
  handleCancel?: (view: EditorView, uploadId: string) => void;
};

function Component(props: NodeViewProps) {
  const { handleRetry, handleCancel } = props.extension.options as Options;
  const uploadId = props.node.attrs["data-upload-id"];
  return (
    <ContentImage
      attrs={props.node.attrs}
      isEditable={props.editor.isEditable}
      isSelected={props.selected && props.editor.isEditable}
      onRetry={() => handleRetry?.(props.view, uploadId)}
      onCancel={() => handleCancel?.(props.view, uploadId)}
      Wrapper={NodeViewWrapper}
    />
  );
}

export const ImageExtended = ImageNode.extend<ImageOptions & Options>({
  addNodeView() {
    return ReactNodeViewRenderer(Component);
  },
  addProseMirrorPlugins() {
    const handleFiles = this.options.handleFiles;
    return [
      // Position tracking plugin - maintains a map of uploadId -> position
      new Plugin({
        key: uploadPositionsKey,
        state: {
          init() {
            return new Map<string, number>();
          },
          apply(_tr, _oldMap, _oldState, newState) {
            // Rebuild position map by scanning the document
            const newMap = new Map<string, number>();

            newState.doc.descendants((node, pos) => {
              const uploadId = node.attrs["data-upload-id"];

              if (uploadId) {
                newMap.set(uploadId, pos);
              }
            });

            return newMap;
          },
        },
      }),

      // Paste handler plugin
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
            void handleFiles?.(view, images);
            return true;
          },
        },
      }),
    ];
  },
});
