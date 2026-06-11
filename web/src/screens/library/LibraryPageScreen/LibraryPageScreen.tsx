"use client";

import { last } from "lodash";
import { useParams } from "next/navigation";
import { parseAsBoolean, parseAsString, useQueryState } from "nuqs";
import { memo } from "react";

import { useNodeGet } from "@/api/openapi-client/nodes";
import { UnreadyBanner } from "@/components/site/Unready";
import { LStack } from "@/styled-system/jsx";

import { Params } from "../library-path";

import { LibraryPageProvider, Props, useLibraryPageContext } from "./Context";
import { EditingDraftWarning } from "./EditingDraftWarning";
import { LibraryPageAutosaveController } from "./LibraryPageAutosaveController";
import { LibraryPageControls } from "./LibraryPageControls";
import { LibraryPageVersionReview } from "./LibraryPageVersionReview";
import { LibraryPageBlocks } from "./blocks/LibraryPageBlocks";
import { LibraryPageEditProvider, useEditState } from "./useEditState";

export function LibraryPageScreen(props: Props) {
  const { slug } = useParams<Params>();
  const [editing] = useQueryState("edit", {
    ...parseAsBoolean,
    defaultValue: false,
    clearOnDefault: true,
  });

  // NOTE: Will fail if slug changes during edit mode.
  const targetSlug = last(slug) ?? props.node.slug;

  const { data, error } = useNodeGet(targetSlug, undefined, {
    swr: {
      fallbackData: props.node,
      // NOTE: We disable all of useSWR's revalidation features while editing
      // in order to not overwrite the author's current state. This isn't great
      // because it means multiple editors will overwrite each others' work. But
      // it's the best we can do without implementing a full sync engine. Yikes.
      enabled: !editing,
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
      <LibraryPageEditProvider>
        <LibraryPageAutosaveController />
        <LibraryPage />
      </LibraryPageEditProvider>
    </LibraryPageProvider>
  );
});
LibraryPageForm.displayName = "LibraryPageForm";

export function LibraryPage() {
  const { initialNode, revalidate } = useLibraryPageContext();
  const [versionID] = useQueryState("version", {
    ...parseAsString,
    clearOnDefault: true,
  });

  const { editMode } = useEditState();

  const editingDraft = editMode === "proposal";

  return (
    <LStack h="full" gap="3" alignItems="start">
      <LibraryPageControls />
      {editingDraft && <EditingDraftWarning />}
      {versionID ? (
        <LibraryPageVersionReview
          node={initialNode}
          versionID={versionID}
          onApplied={revalidate}
        />
      ) : (
        <LibraryPageBlocks />
      )}
    </LStack>
  );
}
