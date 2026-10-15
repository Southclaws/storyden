import { mergeAttributes } from "@tiptap/core";
import Image from "@tiptap/extension-image";

// Persisted content uses this tag; it must remain stable.
const COMPONENT_NAME = "img";

export const ImageNode = Image.extend({
  content: "inline*",
  addOptions() {
    return {
      inline: false,
      allowBase64: false,
      HTMLAttributes: {},
      resize: false,
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
      "data-diff": {
        default: null,
        parseHTML: (element) => element.getAttribute("data-diff"),
        renderHTML: (attributes) => {
          if (!attributes["data-diff"]) {
            return {};
          }
          return {
            "data-diff": attributes["data-diff"],
          };
        },
      },
    };
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
});
