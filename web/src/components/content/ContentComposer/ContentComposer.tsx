"use client";

import { Suspense, lazy } from "react";

import { useHydrated } from "@/utils/useHydrated";

import { ContentRenderer } from "../ContentRenderer/ContentRenderer";
import type { ContentComposerProps } from "../composer-props";

const ContentEditor = lazy(() =>
  import("./ContentEditor").then((module) => ({
    default: module.ContentEditor,
  })),
);

export function ContentComposer(props: ContentComposerProps) {
  const isHydrated = useHydrated();

  if (props.disabled || !isHydrated) {
    return <ContentRenderer {...props} />;
  }

  return (
    <Suspense fallback={<ContentRenderer {...props} />}>
      <ContentEditor {...props} />
    </Suspense>
  );
}
