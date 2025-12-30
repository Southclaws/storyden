import { useEffect, useRef, useState } from "react";

import { handle } from "@/api/client";
import { tagList } from "@/api/openapi-client/tags";
import { InstanceCapability, TagNameList } from "@/api/openapi-schema";
import { CombotagsHandle, CombotagsItem } from "@/components/ui/combotags";
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
  const ref = useRef<CombotagsHandle>(null);
  const [loadingTags, setLoadingTags] = useState(false);

  const currentTags = tags.map((t) => t.name);

  useEffect(() => {
    if (ref.current) {
      ref.current.setValue(currentTags);
    }
  }, [currentTags]);

  async function handleQuery(q: string): Promise<CombotagsItem[]> {
    const tags =
      (await handle(async () => {
        const { tags } = await tagList({ q });
        return tags.map((t) => t.name);
      })) ?? [];

    const filtered = tags.filter((t) => !currentTags.includes(t));

    return filtered.map((name) => ({ id: name, label: name }));
  }

  async function handleSuggestTags() {
    await handle(
      async () => {
        if (!content) {
          throw new Error("Content is required to suggest tags.");
        }

        setLoadingTags(true);
        const tags = await suggestTags(nodeID, content);

        if (!tags) {
          throw new Error(
            "No tags could be suggested for this page. This may be due to the content being too short.",
          );
        }

        ref.current?.append(tags);
      },
      {
        async cleanup() {
          setLoadingTags(false);
        },
      },
    );
  }

  async function handleChange(tags: string[]) {
    setTags(tags);
  }

  return {
    ref,
    currentTags,
    isSuggestEnabled,
    loadingTags,
    handleQuery,
    handleSuggestTags,
    handleChange,
  };
}
