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
import { Button } from "@/components/ui/button";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { LinkButton } from "@/components/ui/link-button";
import { css } from "@/styled-system/css";
import { Center, LStack, styled } from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";

const TAG = "div";

export type LinkPreviewAttributes = {
  href: string;
  "data-display": "card";
};

function LinkPreviewComponent(props: NodeViewProps) {
  const isSelected = props.selected;
  const href = props.node.attrs["href"] as string;
  const isEditable = props.editor.isEditable;

  const { data, error, isMutating, trigger } = useLinkCreate();

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
        width: "full",
        outlineWidth: isSelected ? "medium" : "none",
        outlineStyle: "solid",
        outlineColor: isSelected ? "blue.a6" : "transparent",
        borderRadius: "lg",
        userSelect: isEditable ? "none" : "auto",
        // subtle saturation bump, combined with...
        saturate: isSelected && !isMutating ? "150%" : "100%",
        filter: "auto",
        // background mix with subtle selection colour
        background: isSelected && !isMutating ? "blue.5" : "transparent",
        mixBlendMode: isSelected && !isMutating ? "screen" : "normal",
      })}
    >
      <div data-no-typography>
        {!data ? (
          <Center w="full" minH="12">
            {error ? (
              isEditable ? (
                <styled.div
                  position="absolute"
                  inset="0"
                  display="flex"
                  flexDirection="column"
                  alignItems="center"
                  justifyContent="center"
                  backgroundColor="bg.error"
                  borderRadius="lg"
                  padding="2"
                  height="min"
                  gap="2"
                  userSelect="none"
                  contentEditable={false}
                  role="alert"
                  aria-live="polite"
                >
                  <styled.p
                    fontSize="sm"
                    color="fg.error"
                    fontWeight="medium"
                    maxW="prose"
                  >
                    Link preview failed: {deriveError(error)}
                  </styled.p>
                  <Button
                    type="button"
                    size="xs"
                    variant="subtle"
                    onClick={() => trigger({ url: href })}
                    loading={isMutating}
                  >
                    Retry
                  </Button>
                </styled.div>
              ) : (
                <LStack w="full" gap="1" userSelect="none">
                  <LinkButton size="xs" variant="subtle" href={href}>
                    {href}
                  </LinkButton>
                  <styled.p fontSize="xs" color="fg.muted">
                    <WarningIcon w="3" display="inline" />
                    &nbsp;<span>Link preview failed to load</span>
                  </styled.p>
                </LStack>
              )
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
