import { lazy } from "react";

import { useSession } from "@/auth";

import type { ContentComposerProps } from "../composer-props";

const ContentComposerRich = lazy(() =>
  import("../ContentComposerRich/ContentComposerRich").then((module) => ({
    default: module.ContentComposerRich,
  })),
);

const ContentComposerMarkdown = lazy(() =>
  import("../ContentComposerMarkdown/ContentComposerMarkdown").then(
    (module) => ({ default: module.ContentComposerMarkdown }),
  ),
);

export function ContentEditor(props: ContentComposerProps) {
  const session = useSession();

  return session?.meta.editor.mode === "markdown" ? (
    <ContentComposerMarkdown {...props} />
  ) : (
    <ContentComposerRich {...props} />
  );
}
