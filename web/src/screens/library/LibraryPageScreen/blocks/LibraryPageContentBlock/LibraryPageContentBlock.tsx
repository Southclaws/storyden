import { useState } from "react";

import { ContentComposer } from "@/components/content/ContentComposer/ContentComposer";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useAssets } from "../../useAssets";
import { useEditState } from "../../useEditState";

import { useLibraryContentEvent } from "./events";

export function useLibraryPageContentBlock() {
  const { editing, editorSourceKey, isProposalEditing } = useEditState();
  const { nodeID, store } = useLibraryPageContext();
  const { handleAssetUpload } = useAssets(nodeID);
  const content = useWatch((s) => s.draft.content);
  const { setContent } = store.getState();
  const [contentRevision, setContentRevision] = useState(0);

  function handleChange(value: string) {
    setContent(value);
  }

  useLibraryContentEvent(
    "library-content:update-generated",
    (newContent: string) => {
      setContent(newContent);
      setContentRevision((previous) => previous + 1);
    },
  );

  return {
    editing,
    editorSourceKey,
    isProposalEditing,
    contentRevision,
    content,
    handleAssetUpload,
    handleChange,
  };
}

export function LibraryPageContentBlock() {
  const {
    editing,
    editorSourceKey,
    isProposalEditing,
    contentRevision,
    content,
    handleAssetUpload,
    handleChange,
  } = useLibraryPageContentBlock();

  return (
    <ContentComposer
      key={`content:${editorSourceKey}:${contentRevision}`}
      onChange={handleChange}
      disabled={!editing}
      onAssetUpload={isProposalEditing ? undefined : handleAssetUpload}
      initialValue={content}
    />
  );
}
