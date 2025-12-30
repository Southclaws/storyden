import { useState } from "react";

import { handle } from "@/api/client";
import { tagList } from "@/api/openapi-client/tags";
import { InstanceCapability } from "@/api/openapi-schema";
import { MultiSelectPickerItem } from "@/components/ui/MultiSelectPicker";
import { useLibraryMutation } from "@/lib/library/library";
import { useCapability } from "@/lib/settings/capabilities";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";

export function useLibraryPageTagsBlockEditing() {
  const { nodeID, store } = useLibraryPageContext();
  const { setTags } = store.getState();
  const tags = useWatch((s) => s.draft.tags);
  const content = useWatch((s) => s.draft.content);

  const isSuggestEnabled = useCapability(InstanceCapability.gen_ai);

  const { suggestTags } = useLibraryMutation();
  const [loadingTags, setLoadingTags] = useState(false);
  const [queryResults, setQueryResults] = useState<MultiSelectPickerItem[]>(
    [],
  );

  const currentTagItems: MultiSelectPickerItem[] = tags.map((t) => ({
    label: t.name,
    value: t.name,
  }));

  function handleQuery(q: string) {
    handle(async () => {
      const { tags } = await tagList({ q });
      const currentTagNames = currentTagItems.map((t) => t.value);
      const filtered = tags.filter((t) => !currentTagNames.includes(t.name));
      setQueryResults(
        filtered.map((t) => ({
          label: t.name,
          value: t.name,
        })),
      );
    });
  }

  async function handleSuggestTags() {
    await handle(
      async () => {
        if (!content) {
          throw new Error("Content is required to suggest tags.");
        }

        setLoadingTags(true);
        const suggestedTags = await suggestTags(nodeID, content);

        if (!suggestedTags) {
          throw new Error(
            "No tags could be suggested for this page. This may be due to the content being too short.",
          );
        }

        const currentTagNames = currentTagItems.map((t) => t.value);
        const newTags = [...currentTagNames, ...suggestedTags];
        setTags(newTags);
      },
      {
        async cleanup() {
          setLoadingTags(false);
        },
      },
    );
  }

  async function handleChange(items: MultiSelectPickerItem[]) {
    const tagNames = items.map((item) => item.value);
    setTags(tagNames);
  }

  return {
    currentTagItems,
    queryResults,
    isSuggestEnabled,
    loadingTags,
    handleQuery,
    handleSuggestTags,
    handleChange,
  };
}
