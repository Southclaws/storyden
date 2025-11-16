"use client";

import { useSession } from "@/auth";
import { CardBox } from "@/styled-system/jsx";
import { center } from "@/styled-system/patterns";

import { ContentComposerMarkdown } from "../ContentComposerMarkdown/ContentComposerMarkdown";
import { ContentComposerRich } from "../ContentComposerRich/ContentComposerRich";
import { ContentComposerProps } from "../composer-props";

export function ContentComposer(props: ContentComposerProps) {
  const session = useSession();

  const editorMode = session?.meta.editor.mode;

  switch (editorMode) {
    case "richtext":
      return <ContentComposerRich {...props} />;

    case "markdown":
      return <ContentComposerMarkdown {...props} />;

    default:
      return <ContentComposerRich {...props} />;
  }
}
