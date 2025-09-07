import { useEffect, useState } from "react";

import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useAssets } from "../../useAssets";
import { useEditState } from "../../useEditState";

import { useLibraryContentEvent } from "./events";

export function useLibraryPageContentBlock() {
  const content = useWatch((s) => s.draft.content);

  const [generatedContent, setGeneratedContent] = useState<string | undefined>(
    undefined,
  );

  function handleResetGeneratedContent() {
    setGeneratedContent(undefined);
  }

  // NOTE: We use events for external updates to content in order to side-step
  // React's state management and keep the re-render logic really simple as this
  // component is technically uncontrolled but flips into a controlled state on
  // demand when an external update occurs such as from AI generated content.
  useLibraryContentEvent(
    "library-content:update-generated",
    (newContent: string) => {
      setGeneratedContent(newContent);
    },
  );

  useEffect(() => {
    if (content && generatedContent) {
      // if the content field changes, and it was previously using a controlled
      // value, reset this controlled value to move it back to uncontrolled.
      setGeneratedContent(undefined);
    }
  }, [content, generatedContent]);

  return {
    handleResetGeneratedContent,
    generatedContent,
    content,
  };
}

export function LibraryPageContentBlock() {
  const { nodeID, store } = useLibraryPageContext();
  const { setContent } = store.getState();
  const { editing } = useEditState();
  const { handleAssetUpload } = useAssets(nodeID);
  const { content, generatedContent } = useLibraryPageContentBlock();

  function handleChange(value: string) {
    setContent(value);
  }

  return (
    <ContentComposer
      onChange={handleChange}
      disabled={!editing}
      onAssetUpload={handleAssetUpload}
      initialValue={content}
      value={generatedContent}
    />
  );
}
