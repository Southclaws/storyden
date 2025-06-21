import slugify from "@sindresorhus/slugify";
import { useEffect, useState } from "react";

import { handle } from "@/api/client";
import { InstanceCapability, NodeWithChildren } from "@/api/openapi-schema";
import { useLibraryMutation } from "@/lib/library/library";
import { useCapability } from "@/lib/settings/capabilities";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";

export function useLibraryPageTitleBlock() {
  const { store } = useLibraryPageContext();
  const { draft, setName } = store.getState();

  const { suggestTitle } = useLibraryMutation(draft);
  const [value, setValue] = useState<string | undefined>(undefined);
  const [isLoading, setLoading] = useState(false);
  const isTitleSuggestEnabled = useCapability(InstanceCapability.gen_ai);

  const defaultValue = store.getInitialState().draft.name;
  const name = useWatch((s) => s.draft.name);
  const content = useWatch((s) => s.draft.content);

  // TODO: Figure out an isDirty approach
  // Update the slug with a slugified version of the name if it's not dirty.
  // useEffect(() => {
  //   if (!form.getFieldState("slug").isDirty) {
  //     const autoSlug = slugify(name);
  //     form.setValue("slug", autoSlug);
  //   }
  // }, [form, name]);

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

        const title = await suggestTitle(draft.id, content);
        if (!title) {
          throw new Error("No title could be suggested for this content.");
        }

        setName(title);
        setValue(title);
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
    value,
    isLoading,
    handleReset,
    handleSuggest,
    handleChange,
  };
}
