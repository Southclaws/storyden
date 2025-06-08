"use client";

import { useNodeGet } from "@/api/openapi-client/nodes";
import { UnreadyBanner } from "@/components/site/Unready";
import { LStack } from "@/styled-system/jsx";

import { LibraryPageProvider, Props } from "./Context";
import { LibraryPageControls } from "./LibraryPageControls";
import { LibraryPageBlocks } from "./blocks/LibraryPageBlocks";

export function LibraryPageScreen(props: Props) {
  const { data, error } = useNodeGet(props.node.slug, undefined, {
    swr: { fallbackData: props.node },
  });
  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  return <LibraryPageForm node={data} />;
}

function LibraryPageForm(props: Props) {
  return (
    <LibraryPageProvider node={props.node}>
      <LibraryPage />
    </LibraryPageProvider>
  );
}

export function LibraryPage() {
  return (
    <LStack h="full" gap="3" pl="3" alignItems="start">
      <LibraryPageControls />
      <LibraryPageBlocks />
    </LStack>
  );
}
