import type { PropsWithChildren } from "react";
import "server-only";

import { ContentRendererBody } from "./ContentRendererBody";
import { RenderedContentProvider } from "./RenderedContentContext";

// Pass rendered React nodes across the RSC boundary, without shipping the
// TipTap schema or parser to readers. Match by source so client updates never
// accidentally reuse stale server output.
export function ServerRenderedContent({
  values,
  children,
}: PropsWithChildren<{ values: (string | undefined)[] }>) {
  const content = Object.fromEntries(
    [...new Set(values)]
      .filter((html): html is string => Boolean(html))
      .map((html) => [html, <ContentRendererBody key={html} html={html} />]),
  );
  return (
    <RenderedContentProvider content={content}>
      {children}
    </RenderedContentProvider>
  );
}
