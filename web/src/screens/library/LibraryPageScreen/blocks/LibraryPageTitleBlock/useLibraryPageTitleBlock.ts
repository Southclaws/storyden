import { useRef, useState } from "react";

import { handle } from "@/api/client";
import { InstanceCapability } from "@/api/openapi-schema";
import { useLibraryMutation } from "@/lib/library/library";
import { useCapability } from "@/lib/settings/capabilities";

import { useLibraryPageContext } from "../../Context";
import { useEditState } from "../../useEditState";
import { useWatch } from "../../store";

export function useLibraryPageTitleBlock() {
  const { store } = useLibraryPageContext();
  const { draft, setName } = store.getState();
  const { editorSourceKey } = useEditState();
  const [titleResetVersion, setTitleResetVersion] = useState(0);

  const { suggestTitle } = useLibraryMutation(draft);
  const [isLoading, setLoading] = useState(false);
  const isTitleSuggestEnabled = useCapability(InstanceCapability.gen_ai);

  const value = useWatch((s) => s.draft.name);
  const titleInputKey = `${editorSourceKey}:${titleResetVersion}`;
  const defaultValueRef = useRef<{ key: string; value: string } | null>(null);
  if (defaultValueRef.current?.key !== titleInputKey) {
    defaultValueRef.current = { key: titleInputKey, value };
  }

  const defaultValue = defaultValueRef.current.value;
  const content = useWatch((s) => s.draft.content);

  async function handleSuggest() {
    if (!isTitleSuggestEnabled) {
      return;
    }

    await handle(
      async () => {
        if (!content) {
          throw new Error("Content is required to suggest a title.");
        }

        setLoading(true);

        const title = await suggestTitle(draft.id, content);
        if (!title) {
          throw new Error("No title could be suggested for this content.");
        }

        setName(title);
        setTitleResetVersion((version) => version + 1);
      },
      {
        cleanup: async () => setLoading(false),
      },
    );
  }

  function handleChange(v: string) {
    setName(v);
  }

  return {
    defaultValue,
    isTitleSuggestEnabled,
    titleInputKey,
    isLoading,
    handleSuggest,
    handleChange,
  };
}
