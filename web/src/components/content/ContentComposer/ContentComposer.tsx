"use client";

import { useSession } from "@/auth";
import { CardBox } from "@/styled-system/jsx";
import { center } from "@/styled-system/patterns";

import { ContentComposerMarkdown } from "../ContentComposerMarkdown/ContentComposerMarkdown";
import { ContentComposerRich } from "../ContentComposerRich/ContentComposerRich";
import { ContentComposerProps } from "../composer-props";

export function ContentComposer(props: ContentComposerProps) {
  const session = useSession();

  if (!session) {
    return (
      <CardBox className={center()} p="8" color="fg.muted">
        You must be signed in to compose content
      </CardBox>
    );
  }

  const editorMode = session.meta.editor.mode;

  switch (editorMode) {
    case "richtext":
      return <ContentComposerRich {...props} />;

    case "markdown":
      return <ContentComposerMarkdown {...props} />;

    default:
      return <ContentComposerRich {...props} />;
  }
}
