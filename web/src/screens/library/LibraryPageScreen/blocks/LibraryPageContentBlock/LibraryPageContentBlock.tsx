import { useEffect, useState } from "react";

import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useAssets } from "../../useAssets";
import { useEditState } from "../../useEditState";

export function useLibraryPageContentBlock() {
  const content = useWatch((s) => s.draft.content);

  const [generatedContent, setGeneratedContent] = useState<string | undefined>(
    undefined,
  );

  function handleResetGeneratedContent() {
    setGeneratedContent(undefined);
  }

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
  };
}

export function LibraryPageContentBlock() {
  const { store, currentNode } = useLibraryPageContext();
  const { setContent } = store.getState();
  const { editing } = useEditState();
  const { handleAssetUpload } = useAssets(currentNode);
  const { generatedContent } = useLibraryPageContentBlock();

  function handleChange(value: string) {
    setContent(value);
  }

  return (
    <ContentComposer
      onChange={handleChange}
      disabled={!editing}
      onAssetUpload={handleAssetUpload}
      initialValue={currentNode.content}
      value={generatedContent}
    />
  );
}
