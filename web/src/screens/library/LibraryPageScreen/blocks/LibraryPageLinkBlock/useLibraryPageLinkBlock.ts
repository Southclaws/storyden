"use client";

import { debounce } from "lodash";
import { useCallback, useRef, useState } from "react";
import { toast } from "sonner";

import { handle } from "@/api/client";
import { linkCreate } from "@/api/openapi-client/links";
import { LinkReference } from "@/api/openapi-schema";
import { useLibraryMutation } from "@/lib/library/library";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";

export function useLibraryPageLinkBlock() {
  const { currentNode, store } = useLibraryPageContext();
  const { setLink, setName, setTags } = store.getState();
  const tags = useWatch((s) => s.draft.tags);

  const [inputValue, setInputValue] = useState("");
  const [resolvedLink, setResolvedLink] = useState<
    LinkReference | null | undefined
  >(null);
  const [isImporting, setIsImporting] = useState(false);

  const { revalidate, importFromLink } = useLibraryMutation(currentNode);

  async function handleImportFromLink(link: LinkReference) {
    await handle(
      async () => {
        const { title_suggestion, tag_suggestions, content_suggestion } =
          await importFromLink(currentNode.slug, link.url);

        // TODO: Expose this from suggestion hooks
        // setGeneratedTitle(title_suggestion);
        // setGeneratedContent(content_suggestion);

        if (title_suggestion) {
          setName(title_suggestion);
        }

        setTags([...tags.map((t) => t.name), ...tag_suggestions]);
        // form.setValue("content", content_suggestion);
      },
      {
        cleanup: async () => await revalidate(),
      },
    );
  }

  const debouncedAttemptResolve = useRef(
    debounce(async (s: string) => {
      if (s === "") {
        setResolvedLink(null);
        return;
      }

      let u: URL;
      try {
        u = new URL(s);
      } catch (_) {
        // do nothing for invalid URL.
        setResolvedLink(null);
        return;
      }

      const url = u.toString();

      setLink(url);

      await handle(async () => {
        setResolvedLink(undefined);
        const link = await linkCreate({ url });
        setResolvedLink(link);
      });
    }, 500),
  ).current;

  const handleInputValueChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      const value = event.target.value;
      setInputValue(value);
      debouncedAttemptResolve(value);
    },
    [debouncedAttemptResolve],
  );

  async function handleImport() {
    if (!resolvedLink) {
      toast.error("No link available to import.");
      return;
    }

    setIsImporting(true);
    await handleImportFromLink(resolvedLink);
    setIsImporting(false);
  }

  return {
    data: {
      inputValue,
      resolvedLink,
      isImporting,
    },
    handlers: {
      handleInputValueChange,
      handleImport,
    },
  };
}
