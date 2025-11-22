import Image, { ImageOptions } from "@tiptap/extension-image";
import {
  NodeViewProps,
  NodeViewWrapper,
  ReactNodeViewRenderer,
  mergeAttributes,
} from "@tiptap/react";
import { Plugin, PluginKey } from "prosemirror-state";
import { EditorView } from "prosemirror-view";

import { Asset } from "src/api/openapi-schema";

import { Button } from "@/components/ui/button";
import { ProgressCircle } from "@/components/ui/progress";
import { css } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";

// NOTE: This is the name of the component that will be used in the HTML.
// It cannot be changed.
const COMPONENT_NAME = "img";

// NOTE: This plugin key is used to store the upload positions in ProseMirror
// state to remove the need to scan the whole document for positions.
export const uploadPositionsKey = new PluginKey<Map<string, number>>(
  "uploadPositions",
);

type Options = {
  handleFiles: (view: EditorView, files: File[]) => Promise<Asset[]>;
  handleRetry: (view: EditorView, uploadId: string) => void;
  handleCancel: (view: EditorView, uploadId: string) => void;
};

function Component(props: NodeViewProps) {
  const isUploading = props.node.attrs["data-uploading"] === "true";
  const uploadError = props.node.attrs["data-upload-error"];
  const uploadId = props.node.attrs["data-upload-id"];
  const uploadProgress = props.node.attrs["data-upload-progress"];
  const progressPercent = uploadProgress ? parseInt(uploadProgress, 10) : 0;

  // Access extension options to get retry/cancel handlers
  const { handleRetry, handleCancel } = props.extension.options as Options;

  return (
    <NodeViewWrapper
      className={css({
        position: "relative",
        display: "inline-block",
        cursor: "pointer",
      })}
    >
      <styled.img
        borderRadius="md"
        opacity={isUploading ? "5" : "full"}
        transition="all"
        {...props.node.attrs}
      />
      {isUploading && (
        <styled.div
          position="absolute"
          top="0"
          left="0"
          width="full"
          height="full"
          display="flex"
          flexDirection="column"
          alignItems="center"
          justifyContent="center"
          pointerEvents="none"
          gap="3"
          padding="4"
        >
          <ProgressCircle value={progressPercent} size="md" />
        </styled.div>
      )}
      {uploadError && (
        <styled.div
          position="absolute"
          inset="0"
          display="flex"
          flexDirection="column"
          alignItems="center"
          justifyContent="center"
          backgroundColor="bg.error"
          opacity="9"
          borderRadius="md"
          padding="3"
          gap="2"
          userSelect="none"
          contentEditable={false}
        >
          <styled.p fontSize="sm" color="fg.error" fontWeight="medium">
            Upload failed
          </styled.p>
          <styled.div display="flex" gap="2">
            <Button
              type="button"
              size="xs"
              variant="outline"
              onClick={() => handleRetry(props.view, uploadId)}
            >
              Retry
            </Button>
            <Button
              type="button"
              size="xs"
              variant="ghost"
              onClick={() => handleCancel(props.view, uploadId)}
            >
              Remove
            </Button>
          </styled.div>
        </styled.div>
      )}
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
  addAttributes() {
    return {
      ...this.parent?.(),
      "data-upload-id": {
        default: null,
        parseHTML: (element) => element.getAttribute("data-upload-id"),
        renderHTML: (attributes) => {
          if (!attributes["data-upload-id"]) {
            return {};
          }
          return {
            "data-upload-id": attributes["data-upload-id"],
          };
        },
      },
      "data-uploading": {
        default: null,
        parseHTML: (element) => element.getAttribute("data-uploading"),
        renderHTML: (attributes) => {
          if (!attributes["data-uploading"]) {
            return {};
          }
          return {
            "data-uploading": attributes["data-uploading"],
          };
        },
      },
      "data-upload-error": {
        default: null,
        parseHTML: (element) => element.getAttribute("data-upload-error"),
        renderHTML: (attributes) => {
          if (!attributes["data-upload-error"]) {
            return {};
          }
          return {
            "data-upload-error": attributes["data-upload-error"],
          };
        },
      },
      "data-upload-progress": {
        default: null,
        parseHTML: (element) => element.getAttribute("data-upload-progress"),
        renderHTML: (attributes) => {
          if (!attributes["data-upload-progress"]) {
            return {};
          }
          return {
            "data-upload-progress": attributes["data-upload-progress"],
          };
        },
      },
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
              // const isUploading = node.attrs["data-uploading"] === "true";

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
            handleFiles(view, images);
            return true;
          },
        },
      }),
    ];
  },
});
