import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { KeyedMutator } from "swr";
import { z } from "zod";

import { linkCreate } from "src/api/openapi-client/links";
import { LinkListResult, LinkReference } from "src/api/openapi-schema";

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

  const { q } = form.watch();

  const url = getURL(q);

  function resetIndexingState() {
    setTimeout(() => {
      setIndexing({
        state: "not-indexing",
      });
    }, indexingResetTimeout);
  }

  useEffect(() => {
    if (!url) {
      setIndexing({
        state: "not-indexing",
      });

      return;
    }

    (async () => {
      setIndexing({
        state: "indexing",
        url,
      });

      try {
        const link = await linkCreate({
          url,
        });

        setIndexing({
          state: "indexed",
          link,
        });
      } catch (_) {
        setIndexing({
          state: "error",
          error: "Failed to index the provided link.",
        });

        resetIndexingState();
        return;
      }
    })();
  }, [url]);

  const handlePage = async (page: number) => {
    router.push(`/links?q=${q}&page=${page}`);
  };

  const handleSubmission = form.handleSubmit(async (payload) => {
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

function getURL(s: string) {
  try {
    if (s.length < 5) {
      return undefined;
    }

    if (s.includes(" ")) {
      return undefined;
    }

    const parsed = new URL(s);

    return parsed.toString();
  } catch (_) {
    return undefined;
  }
}
