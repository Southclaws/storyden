"use client";

import { debounce } from "lodash";
import { useCallback, useEffect, useRef, useState } from "react";
import { toast } from "sonner";

import { handle } from "@/api/client";
import { linkCreate } from "@/api/openapi-client/links";
import { LinkReference } from "@/api/openapi-schema";
import { useLibraryMutation } from "@/lib/library/library";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";

export function useLibraryPageLinkBlock() {
  const { nodeID, store } = useLibraryPageContext();
  const { setLink, setName, setTags } = store.getState();
  const tags = useWatch((s) => s.draft.tags);
  const link = useWatch((s) => s.draft.link);

  const defaultLinkURL = link?.url || "";
  const [inputValue, setInputValue] = useState(defaultLinkURL);
  const [resolvedLink, setResolvedLink] = useState<
    LinkReference | null | undefined
  >(null);
  const [isImporting, setIsImporting] = useState(false);

  const { revalidate, importFromLink } = useLibraryMutation();

  async function handleImportFromLink(link: LinkReference) {
    await handle(
      async () => {
        const { title_suggestion, tag_suggestions, content_suggestion } =
          await importFromLink(nodeID, link.url);

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

      await handle(async () => {
        setResolvedLink(undefined);
        const link = await linkCreate({ url });
        setResolvedLink(link);
        setLink(link);
      });
    }, 500),
  ).current;

  // Cleanup debounced function on component unmount
  useEffect(() => {
    return () => {
      debouncedAttemptResolve.cancel();
    };
  }, [debouncedAttemptResolve]);

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
      defaultLinkURL,
      isImporting,
    },
    handlers: {
      handleInputValueChange,
      handleImport,
    },
  };
}
