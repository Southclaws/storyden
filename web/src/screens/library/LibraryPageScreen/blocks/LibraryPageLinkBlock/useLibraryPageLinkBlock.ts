"use client";

import { debounce } from "lodash";
import { useCallback, useEffect, useRef, useState } from "react";
import { toast } from "sonner";

import { handle } from "@/api/client";
import { linkCreate } from "@/api/openapi-client/links";
import { LinkReference } from "@/api/openapi-schema";
import { isContentEmpty } from "@/lib/content/content";
import { useLibraryMutation } from "@/lib/library/library";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useEmitLibraryContentEvent } from "../LibraryPageContentBlock/events";
import { useEmitLibraryCoverEvent } from "../LibraryPageCoverBlock/events";

export function useLibraryPageLinkBlock() {
  const { nodeID, store } = useLibraryPageContext();
  const { setLink, removeLink, setName, setTags } = store.getState();
  const tags = useWatch((s) => s.draft.tags);
  const link = useWatch((s) => s.draft.link);
  const content = useWatch((s) => s.draft.content);
  const primaryImage = useWatch((s) => s.draft.primary_image);
  const emitContentEvent = useEmitLibraryContentEvent();
  const emitCoverEvent = useEmitLibraryCoverEvent();

  const defaultLinkURL = link?.url || "";
  const [inputValue, setInputValue] = useState(defaultLinkURL);
  const [resolvedLink, setResolvedLink] = useState<
    LinkReference | null | undefined
  >(link);
  const [isImporting, setIsImporting] = useState(false);

  const { revalidate, importFromLink } = useLibraryMutation();

  const isEmpty = isContentEmpty(content);

  async function handleImportFromLink(link: LinkReference) {
    await handle(
      async () => {
        const {
          title_suggestion,
          tag_suggestions,
          content_suggestion,
          primary_image,
        } = await importFromLink(nodeID, link.url);

        // TODO: Is this actually useful?
        if (title_suggestion) {
          setName(title_suggestion);
        }

        if (tag_suggestions.length > 0) {
          const existingNames = tags.map((t) => t.name);
          const newTags = [...existingNames];

          for (const suggestion of tag_suggestions) {
            if (!existingNames.includes(suggestion)) {
              newTags.push(suggestion);
            }
          }

          setTags(newTags);
        }

        // NOTE: Only apply AI suggestion if the content is empty to avoid
        // overwriting member-written content.
        // TODO: Implement a better UI/UX for this using a kind of AI-proposal.
        if (content_suggestion && isEmpty) {
          emitContentEvent(
            "library-content:update-generated",
            content_suggestion,
          );
        }

        // Set cover image from imported opengraph image if it doesn't have one.
        if (primary_image && primaryImage === undefined) {
          emitCoverEvent("library-cover:update-from-asset", primary_image);
        }
      },
      {
        cleanup: async () => await revalidate(),
      },
    );
  }

  const debouncedAttemptResolve = useRef(
    debounce(async (s: string) => {
      if (s === "") {
        setResolvedLink(undefined);
        removeLink();
        return;
      }

      let u: URL;
      try {
        u = new URL(s);
      } catch (_) {
        // do nothing for invalid URL.
        setResolvedLink(undefined);
        return;
      }

      const url = u.toString();

      await handle(async () => {
        setResolvedLink(null);
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
      const value = event.target.value.trim();
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
