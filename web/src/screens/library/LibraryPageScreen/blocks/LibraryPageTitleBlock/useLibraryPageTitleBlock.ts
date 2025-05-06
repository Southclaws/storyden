import { useState } from "react";

import { handle } from "@/api/client";
import { InstanceCapability, NodeWithChildren } from "@/api/openapi-schema";
import { useLibraryMutation } from "@/lib/library/library";
import { useCapability } from "@/lib/settings/capabilities";

import { useLibraryPageContext } from "../../Context";

export function useLibraryPageTitleBlock() {
  const { form, node } = useLibraryPageContext();
  const { suggestTitle } = useLibraryMutation(node);
  const [value, setValue] = useState<string | undefined>(undefined);
  const [isLoading, setLoading] = useState(false);
  const isTitleSuggestEnabled = useCapability(InstanceCapability.gen_ai);

  function handleReset() {
    setValue(undefined);
  }

  async function handleSuggest() {
    if (!isTitleSuggestEnabled) {
      return;
    }

    await handle(
      async () => {
        setLoading(true);

        const title = await suggestTitle(node.id);
        if (!title) {
          throw new Error("No title could be suggested for this content.");
        }

        form.setValue("name", title);
        setValue(title);
      },
      {
        cleanup: async () => setLoading(false),
      },
    );
  }

  return {
    isTitleSuggestEnabled,
    value,
    isLoading,
    handleReset,
    handleSuggest,
  };
}
