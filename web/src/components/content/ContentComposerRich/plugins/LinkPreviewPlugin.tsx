import {
  type NodeViewProps,
  NodeViewWrapper,
  ReactNodeViewRenderer,
} from "@tiptap/react";

import { ContentLinkCard } from "../../ContentRenderer/ContentLinkCard";
import { LinkPreviewNode } from "../../ContentRenderer/LinkPreviewNode";

export type LinkPreviewAttributes = {
  href: string;
  "data-display": "card";
  "data-diff"?: "insertion" | "deletion";
};

function LinkPreviewComponent(props: NodeViewProps) {
  return (
    <ContentLinkCard
      href={props.node.attrs["href"]}
      dataDiff={props.node.attrs["data-diff"]}
      isEditable={props.editor.isEditable}
      isSelected={props.selected && props.editor.isEditable}
      Wrapper={NodeViewWrapper}
    />
  );
}

export const LinkPreview = LinkPreviewNode.extend({
  addNodeView() {
    return ReactNodeViewRenderer(LinkPreviewComponent);
  },

  addCommands() {
    return {
      setLinkPreview:
        (attributes: Pick<LinkPreviewAttributes, "href">) =>
        ({ commands }) => {
          return commands.insertContent({
            type: this.name,
            attrs: {
              href: attributes.href,
              "data-display": "card",
            } satisfies LinkPreviewAttributes,
          });
        },
    };
  },
});
