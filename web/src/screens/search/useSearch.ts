"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { parseAsInteger, useQueryState } from "nuqs";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useDatagraphSearch } from "@/api/openapi-client/datagraph";
import {
  DatagraphItemKind,
  DatagraphSearchOKResponse,
} from "@/api/openapi-schema";
import { useSearchQueryState } from "@/components/search/Search/useSearch";
import { DatagraphKindSchema } from "@/lib/datagraph/schema";

export type Props = {
  initialQuery: string;
  initialPage: number;
  initialKind: DatagraphItemKind[];
  initialResults?: DatagraphSearchOKResponse;
};

export const FormSchema = z.object({
  q: z.string().min(1, { message: "Please enter a search term" }),
  kind: z.array(DatagraphKindSchema).optional(),
});
export type Form = z.infer<typeof FormSchema>;

export function useSearchScreen(props: Props) {
  const [query, setQuery] = useSearchQueryState();
  const [page, setPage] = useQueryState("page", {
    ...parseAsInteger,
    defaultValue: 1,
  });

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      q: props.initialQuery,
      kind: props.initialKind,
    },
  });

  const kind = form.watch("kind");

  // NOTE: Because useQueryState does not support proper URL query arrays, we
  // modify the "kind" query parameter array directly using browser APIs.
  useEffect(() => {
    if (kind) {
      const url = new URL(window.location.href);

      url.searchParams.delete("kind");
      kind.forEach((k) => url.searchParams.append("kind", k));

      window.history.replaceState({}, "", url.toString());
    }
  }, [kind]);

  const { data, error, isLoading } = useDatagraphSearch(
    {
      q: query,
      kind: kind,
      page: page.toString(),
    },
    {
      swr: {
        enabled: !!query,
        fallbackData: props.initialResults,
      },
    },
  );

  // NOTE: This is done via a useEffect because we don't want this to be present
  // on a server-render, only for client side search interactions.
  const [isSearchLoading, setLoading] = useState(false);
  useEffect(() => {
    setLoading(isLoading ?? false);
  }, [isLoading]);

  const handleSearch = form.handleSubmit((data) => {
    setQuery(data.q);
  });

  const handleReset = async () => {
    form.reset();
    setQuery(null);
  };

  return {
    form,
    isLoading: isSearchLoading,
    error,
    data: {
      query,
      page: page,
      results: data,
    },
    handlers: {
      handleSearch,
      handleReset,
    },
  };
}
