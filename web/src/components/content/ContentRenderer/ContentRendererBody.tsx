import { generateJSON } from "@tiptap/html";
import { renderToReactElement } from "@tiptap/static-renderer/pm/react";

import { ContentImage } from "./ContentImage";
import { ContentLinkCard } from "./ContentLinkCard";
import { contentExtensions } from "./extensions";

export function ContentRendererBody({ html }: { html: string }) {
  const content = generateJSON(html, contentExtensions);
  return renderToReactElement({
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
        image: ({ node }) => <ContentImage attrs={{ ...node.attrs }} />,
        linkPreview: ({ node }) => (
          <ContentLinkCard
            href={node.attrs["href"]}
            dataDiff={node.attrs["data-diff"]}
          />
        ),
      },
    },
  });
}
