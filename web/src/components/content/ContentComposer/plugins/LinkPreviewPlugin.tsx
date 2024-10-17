import {
  Node,
  NodeViewProps,
  mergeAttributes,
  nodePasteRule,
} from "@tiptap/core";
import { NodeViewWrapper, ReactNodeViewRenderer } from "@tiptap/react";

const COMPONENT_NAME = "a";

export interface LinkPreviewOptions {
  /**
   * Controls if the paste handler for youtube videos should be added.
   * @default true
   * @example false
   */
  addPasteHandler: boolean;

  /**
   * The HTML attributes for a youtube video node.
   * @default {}
   * @example { class: 'foo' }
   */
  HTMLAttributes: Record<string, any>;
}

function Component(props: NodeViewProps) {
  return (
    <NodeViewWrapper>
      <pre>[link preview nodeview]</pre>
    </NodeViewWrapper>
  );
}

export const LinkPreview = Node.create<LinkPreviewOptions>({
  name: "linkPreview",

  addOptions() {
    return {
      addPasteHandler: true,
      HTMLAttributes: {
        "data-link-preview": "true",
      },
    };
  },

  inline() {
    return false;
  },

  group() {
    return "block";
  },

  draggable: true,

  addNodeView() {
    return ReactNodeViewRenderer(Component);
  },

  addAttributes() {
    return {
      href: {
        default: null,
      },
    };
  },

  parseHTML() {
    return [
      {
        tag: "a",
        getAttrs: (dom) => {
          const href = (dom as HTMLElement).getAttribute("data-link-preview");

          console.log("href", href);

          return null;
        },
      },
    ];
  },

  addPasteRules() {
    if (!this.options.addPasteHandler) {
      return [];
    }

    return [
      nodePasteRule({
        find: /https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z]{2,}\b(?:[-a-zA-Z0-9@:%._+~#=?!&/]*)(?:[-a-zA-Z0-9@:%._+~#=?!&/]*)/gi,
        type: this.type,
        getAttributes: (match) => {
          return { src: match.input };
        },
        getContent: (x) => {
          console.log(x);
          return [];
        },
      }),
    ];
  },

  renderHTML({ HTMLAttributes }) {
    return [COMPONENT_NAME, mergeAttributes(HTMLAttributes), 0];
  },
});
