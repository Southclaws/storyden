import slugify from "@sindresorhus/slugify";
import { useEffect, useState } from "react";

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

  const { name, content } = form.watch();

  // Update the slug with a slugified version of the name if it's not dirty.
  useEffect(() => {
    if (!form.getFieldState("slug").isDirty) {
      const autoSlug = slugify(name);
      form.setValue("slug", autoSlug);
    }
  }, [form, name]);

  function handleReset() {
    setValue(undefined);
  }

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

        const title = await suggestTitle(node.id, content);
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
