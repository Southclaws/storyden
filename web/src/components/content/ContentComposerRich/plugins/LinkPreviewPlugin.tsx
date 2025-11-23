import { Node } from "@tiptap/core";
import {
  NodeViewProps,
  NodeViewWrapper,
  ReactNodeViewRenderer,
} from "@tiptap/react";
import { useEffect } from "react";

import { useLinkCreate } from "@/api/openapi-client/links";
import { LinkCard } from "@/components/library/links/LinkCard";
import { Spinner } from "@/components/ui/Spinner";
import { css } from "@/styled-system/css";
import { Center, styled } from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";

const TAG = "div";

export type LinkPreviewAttributes = {
  href: string;
  "data-display": "card";
};

function LinkPreviewComponent(props: NodeViewProps) {
  const href = props.node.attrs["href"] as string;
  const isEditable = props.editor.isEditable;

  const { data, error, trigger } = useLinkCreate();

  useEffect(() => {
    trigger({
      url: href,
    });
  }, [href]);

  return (
    <NodeViewWrapper
      className={css({
        position: "relative",
        display: "inline-block",
        cursor: "pointer",
        width: "full",
      })}
    >
      <div data-no-typography>
        {!data ? (
          <Center w="full" h="12">
            {error ? (
              <>
                (
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
                  <styled.p
                    fontSize="sm"
                    color="fg.error"
                    fontWeight="medium"
                    maxW="prose"
                  >
                    Link preview failed: {deriveError(error)}
                  </styled.p>
                </styled.div>
                )
              </>
            ) : (
              <Spinner />
            )}
          </Center>
        ) : (
          <LinkCard shape="row" link={data} disableAnchors={isEditable} />
        )}
      </div>
    </NodeViewWrapper>
  );
}

export const LinkPreview = Node.create<{}>({
  name: "linkPreview",

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
