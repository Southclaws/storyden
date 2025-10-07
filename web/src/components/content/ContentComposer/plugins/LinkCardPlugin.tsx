import { Node, mergeAttributes } from "@tiptap/core";
import {
  NodeViewProps,
  NodeViewWrapper,
  ReactNodeViewRenderer,
} from "@tiptap/react";

import { Box } from "@/styled-system/jsx";

const TAG = "div";

type LinkCardAttributes = {
  href: string;
  "data-display": "card";
};

function LinkCardComponent(props: NodeViewProps) {
  const href = props.node.attrs["href"] as string;

  return (
    <NodeViewWrapper>
      <Box
        padding="4"
        borderColor="border.default"
        borderRadius="md"
        backgroundColor="bg.muted"
        cursor="pointer"
        style={{ borderWidth: "1px", borderStyle: "solid" }}
      >
        <pre>{href}</pre>
      </Box>
    </NodeViewWrapper>
  );
}

export const LinkCard = Node.create<{}>({
  name: "linkCard",

  group: "block",

  atom: true,

  selectable: true,

  addAttributes() {
    return {
      href: {
        default: null,
        parseHTML: (element) => element.getAttribute("data-href"),
        renderHTML: (href) => {
          if (!href) {
            return {};
          }
          return { "data-href": href };
        },
      },
      "data-display": {
        default: "card",
        parseHTML: (element) => element.getAttribute("data-display"),
        renderHTML: (display) => {
          return { "data-display": display };
        },
      },
    };
  },

  parseHTML() {
    return [
      {
        tag: `${TAG}[data-display="card"]`,
        priority: 100,
      },
    ];
  },

  renderHTML({ node }) {
    const href = node.attrs["href"] || "";
    return [
      TAG,
      {
        "data-href": href,
        "data-display": "card",
        class: "link-card",
      },
    ];
  },

  addNodeView() {
    return ReactNodeViewRenderer(LinkCardComponent);
  },

  addCommands() {
    return {
      setLinkCard:
        (attributes: { href: string }) =>
        ({ commands }: any) => {
          return commands.insertContent({
            type: this.name,
            attrs: {
              href: attributes.href,
              "data-display": "card",
            },
          });
        },
    } as any;
  },
});
