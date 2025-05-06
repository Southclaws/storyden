"use client";

import { useState } from "react";
import { toast } from "sonner";

import { handle } from "@/api/client";
import { linkCreate } from "@/api/openapi-client/links";
import { LinkReference } from "@/api/openapi-schema";
import { useLibraryMutation } from "@/lib/library/library";

import { useLibraryPageContext } from "../../Context";

export function useLibraryPageLinkBlock() {
  const { form, node } = useLibraryPageContext();
  const [link, setLink] = useState<LinkReference | null | undefined>(null);
  const [isImporting, setIsImporting] = useState(false);

  const { revalidate, importFromLink } = useLibraryMutation(node);

  async function handleImportFromLink(link: LinkReference) {
    await handle(
      async () => {
        const { title_suggestion, tag_suggestions, content_suggestion } =
          await importFromLink(node.slug, link.url);

        // TODO: Expose this from suggestion hooks
        // setGeneratedTitle(title_suggestion);
        // setGeneratedContent(content_suggestion);

        if (title_suggestion) {
          // TODO: Expose a setName hook for page context.
          form.setValue("name", title_suggestion);
        }
        // TODO: Expose setTags, setContent hooks for page context.
        form.setValue("tags", tag_suggestions);
        form.setValue("content", content_suggestion);
      },
      {
        cleanup: async () => await revalidate(),
      },
    );
  }

  async function handleURL(s: string) {
    if (s === "") {
      setLink(null);
      return;
    }

    await handle(async () => {
      try {
        setLink(undefined);

        const u = new URL(s);

        await new Promise((resolve) => setTimeout(resolve, 1000));

        const link = await linkCreate({ url: u.toString() });
        setLink(link);
      } catch (_) {
        // do nothing for invalid URL, already handled by parent form logic.
        setLink(null);
      }
    });
  }

  async function handleImport() {
    if (!link) {
      toast.error("No link available to import.");
      return;
    }

    setIsImporting(true);
    await handleImportFromLink(link);
    setIsImporting(false);
  }

  return {
    form,
    data: {
      link,
      isImporting,
    },
    handlers: {
      handleURL,
      handleImport,
    },
  };
}
