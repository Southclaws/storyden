import { useRef, useState } from "react";

import { handle } from "@/api/client";
import { tagList } from "@/api/openapi-client/tags";
import { InstanceCapability, TagNameList } from "@/api/openapi-schema";
import { CombotagsHandle } from "@/components/ui/combotags";
import { useLibraryMutation } from "@/lib/library/library";
import { useCapability } from "@/lib/settings/capabilities";

import { useLibraryPageContext } from "../../Context";

export function useLibraryPageTagsBlockEditing() {
  const { node, form } = useLibraryPageContext();
  const isSuggestEnabled = useCapability(InstanceCapability.gen_ai);

  const { suggestTags } = useLibraryMutation(node);
  const ref = useRef<CombotagsHandle>(null);
  const [loadingTags, setLoadingTags] = useState(false);

  const currentTags = node.tags.map((t) => t.name);

  async function handleQuery(q: string): Promise<TagNameList> {
    const tags =
      (await handle(async () => {
        const { tags } = await tagList({ q });
        return tags.map((t) => t.name);
      })) ?? [];

    const filtered = tags.filter((t) => !currentTags.includes(t));

    return filtered;
  }

  async function handleSuggestTags() {
    await handle(
      async () => {
        setLoadingTags(true);
        const tags = await suggestTags(node.slug);

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

  return {
    ref,
    form,
    currentTags,
    isSuggestEnabled,
    loadingTags,
    handleQuery,
    handleSuggestTags,
  };
}
