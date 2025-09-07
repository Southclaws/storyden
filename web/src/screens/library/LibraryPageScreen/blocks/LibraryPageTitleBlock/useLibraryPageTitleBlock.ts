import slugify from "@sindresorhus/slugify";
import { useEffect, useRef, useState } from "react";

import { handle } from "@/api/client";
import { InstanceCapability } from "@/api/openapi-schema";
import { useLibraryMutation } from "@/lib/library/library";
import { useCapability } from "@/lib/settings/capabilities";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";

export function useLibraryPageTitleBlock() {
  const { store } = useLibraryPageContext();
  const { draft, setName, setSlug } = store.getState();

  const { suggestTitle } = useLibraryMutation(draft);
  const [isLoading, setLoading] = useState(false);
  const isTitleSuggestEnabled = useCapability(InstanceCapability.gen_ai);

  const defaultValue = store.getInitialState().draft.name;
  const value = useWatch((s) => s.draft.name);
  const slug = useWatch((s) => s.draft.slug);
  const content = useWatch((s) => s.draft.content);

  // Track the previous title to detect if slug should be auto-updated
  const previousTitleRef = useRef(value);

  // Auto-update slug when title changes, but only if slug appears to be auto-generated
  useEffect(() => {
    const currentTitle = value || "";
    const previousTitle = previousTitleRef.current || "";

    if (currentTitle !== previousTitle && currentTitle.trim()) {
      const shouldUpdateSlug =
        // Slug is empty
        !slug ||
        // Slug contains "untitled" (auto-generated)
        slug.includes("untitled") ||
        // Slug matches the slugified version of the previous title (not user-customized)
        slug === slugify(previousTitle);

      if (shouldUpdateSlug) {
        const newSlug = slugify(currentTitle);
        if (newSlug !== slug) {
          setSlug(newSlug);
        }
      }
    }

    previousTitleRef.current = currentTitle;
  }, [value, slug, setSlug]);

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
    handleSuggest,
    handleChange,
  };
}
