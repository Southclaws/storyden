import { generateJSON } from "@tiptap/html";
import { renderToReactElement } from "@tiptap/static-renderer/pm/react";

import { css, cx } from "@/styled-system/css";
import { LStack } from "@/styled-system/jsx";

import type { ContentComposerProps } from "../composer-props";

import "./styles.css";

import { ContentImage } from "./ContentImage";
import { ContentLinkCard } from "./ContentLinkCard";
import { contentExtensions } from "./extensions";

export function ContentRenderer({
  initialValue,
  value,
  className,
}: ContentComposerProps) {
  const content = generateJSON(
    value ?? initialValue ?? "<p></p>",
    contentExtensions,
  );

  return (
    <LStack
      containerType="inline-size"
      className={cx(
        "content-composer content-composer--readonly typography",
        className,
      )}
      position="relative"
      w="full"
      gap="1"
      minHeight="8"
    >
      <div className={css({ height: "full", width: "full" })}>
        <div
          className={cx(
            "content-composer__content tiptap ProseMirror content-renderer",
            css({ height: "full", width: "full" }),
          )}
        >
          {renderToReactElement({
            extensions: contentExtensions,
            content,
            options: {
              nodeMapping: {
                paragraph: ({ node, children }) => (
                  <p>
                    {children}
                    {(!node.childCount ||
                      node.lastChild?.type.name === "hardBreak") && (
                      <br className="ProseMirror-trailingBreak" />
                    )}
                  </p>
                ),
                image: ({ node }) => <ContentImage attrs={node.attrs} />,
                linkPreview: ({ node }) => (
                  <ContentLinkCard
                    href={node.attrs["href"]}
                    dataDiff={node.attrs["data-diff"]}
                  />
                ),
              },
            },
          })}
        </div>
      </div>
    </LStack>
  );
}
