"use client";

import { last } from "lodash";
import { useParams } from "next/navigation";
import { memo } from "react";

import { useNodeGet } from "@/api/openapi-client/nodes";
import { UnreadyBanner } from "@/components/site/Unready";
import { LStack } from "@/styled-system/jsx";

import { Params } from "../library-path";

import { LibraryPageProvider, Props } from "./Context";
import { LibraryPageControls } from "./LibraryPageControls";
import { LibraryPageBlocks } from "./blocks/LibraryPageBlocks";

export function LibraryPageScreen(props: Props) {
  const { slug } = useParams<Params>();

  // NOTE: Will fail if slug changes during edit mode.
  const targetSlug = last(slug) ?? props.node.slug;

  const { data, error } = useNodeGet(targetSlug, undefined, {
    swr: {
      fallbackData: props.node,
      revalidateOnFocus: false,
    },
  });
  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  return <LibraryPageForm node={data} childNodes={props.childNodes} />;
}

const LibraryPageForm = memo((props: Props) => {
  return (
    <LibraryPageProvider node={props.node} childNodes={props.childNodes}>
      <LibraryPage />
    </LibraryPageProvider>
  );
});
LibraryPageForm.displayName = "LibraryPageForm";

export function LibraryPage() {
  return (
    <LStack h="full" gap="3" pl="3" alignItems="start">
      <LibraryPageControls />
      <LibraryPageBlocks />
    </LStack>
  );
}
