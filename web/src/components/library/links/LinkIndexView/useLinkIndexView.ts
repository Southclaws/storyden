import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { KeyedMutator } from "swr";

import { linkCreate } from "@/api/openapi-client/links";
import { LinkListResult, LinkReference } from "@/api/openapi-schema";
import { normalizeLink } from "@/lib/link/validation";
import * as z from "@/lib/validation";

export type Props = {
  links: LinkListResult;
  mutate?: KeyedMutator<LinkListResult>;
  query?: string;
  page?: number;
};

export const FormSchema = z.object({
  q: z.string(),
});
export type Form = z.infer<typeof FormSchema>;

const indexingDelay = 500;
const indexingResetTimeout = 5000;

export type IndexingState =
  | {
      state: "not-indexing";
    }
  | {
      state: "indexing";
      url: string;
    }
  | {
      state: "indexed";
      link: LinkReference;
    }
  | {
      state: "error";
      error: string;
    };

const defaultIndexingState = { state: "not-indexing" } satisfies IndexingState;

export function useLinkIndexView(props: Props) {
  const router = useRouter();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      q: props.query,
    },
  });
  const [indexing, setIndexing] = useState<IndexingState>(defaultIndexingState);
  const indexedURL = useRef<string | undefined>(undefined);

  const { q } = form.watch();

  const url = normalizeLink(q);

  function resetIndexingState() {
    setTimeout(() => {
      setIndexing({
        state: "not-indexing",
      });
    }, indexingResetTimeout);
  }

  async function indexLink(url: string) {
    if (indexedURL.current === url) {
      return;
    }
    indexedURL.current = url;

    setIndexing({
      state: "indexing",
      url,
    });

    try {
      const link = await linkCreate({
        url,
      });
      if (indexedURL.current !== url) {
        return;
      }

      setIndexing({
        state: "indexed",
        link,
      });

      await props.mutate?.();
    } catch (_) {
      if (indexedURL.current !== url) {
        return;
      }

      setIndexing({
        state: "error",
        error: "Failed to index the provided link.",
      });

      resetIndexingState();
    }
  }

  useEffect(() => {
    if (indexedURL.current === url) {
      return;
    }

    indexedURL.current = undefined;
    setIndexing(defaultIndexingState);

    if (!url) {
      return;
    }

    const timeout = setTimeout(() => indexLink(url), indexingDelay);
    return () => clearTimeout(timeout);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [url]);

  const handlePage = async (page: number) => {
    router.push(`/links?q=${q}&page=${page}`);
  };

  const handleSubmission = form.handleSubmit(async (payload) => {
    const submittedURL = normalizeLink(payload.q);
    if (submittedURL) {
      await indexLink(submittedURL);
      return;
    }

    router.push(`/links?q=${payload.q}`);
  });

  const handleReset = async () => {
    form.reset();
    router.push("/links");
  };

  const handleMutate = async () => {
    await props.mutate?.();
  };

  return {
    form,
    data: {
      q,
      isURL: url !== undefined,
      links: props.links,
      indexing,
    },
    handlers: {
      handlePage,
      handleSubmission,
      handleReset,
      handleMutate,
    },
  };
}
