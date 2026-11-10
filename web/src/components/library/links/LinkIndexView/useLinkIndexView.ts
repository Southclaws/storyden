import { parseAsInteger, parseAsString, useQueryStates } from "nuqs";
import { FormEvent, useEffect, useRef, useState } from "react";
import { useSWRConfig } from "swr";

import {
  getLinkListKey,
  linkCreate,
  useLinkList,
} from "@/api/openapi-client/links";
import { LinkListResult, LinkReference } from "@/api/openapi-schema";
import { normalizeLink } from "@/lib/link/validation";
import { deriveError } from "@/utils/error";

export type Props = {
  initialResult?: LinkListResult;
  query?: string;
  page?: number;
};

type AddState =
  | { state: "idle" }
  | { state: "adding" }
  | { state: "added" }
  | { state: "error"; error: string };

export function useLinkIndexView(props: Props) {
  const { mutate } = useSWRConfig();
  const [filters, setFilters] = useQueryStates({
    q: parseAsString.withDefault(""),
    page: parseAsInteger.withDefault(1),
  });
  const query = filters.q.trim();
  const page = Math.max(1, filters.page);
  const [search, setSearch] = useState({ query, page });
  const [adding, setAdding] = useState<AddState>({ state: "idle" });
  const requestRef = useRef<{ url: string } | undefined>(undefined);
  const url = normalizeLink(query);
  const [previousURL, setPreviousURL] = useState(url);

  const isInitialSearch =
    search.query === (props.query ?? "").trim() &&
    search.page === Number(props.page ?? 1);
  const {
    data: links,
    error,
    mutate: refresh,
  } = useLinkList(
    { q: search.query, page: search.page.toString() },
    {
      swr: {
        fallbackData: isInitialSearch ? props.initialResult : undefined,
        keepPreviousData: false,
      },
    },
  );

  useEffect(() => {
    const timeout = setTimeout(() => setSearch({ query, page }), 300);
    return () => clearTimeout(timeout);
  }, [query, page]);

  useEffect(() => {
    return () => {
      requestRef.current = undefined;
    };
  }, [url]);

  if (previousURL !== url) {
    setPreviousURL(url);
    setAdding({ state: "idle" });
  }

  const isSearching = search.query !== query || search.page !== page;
  const hasSavedLink = links?.links.some((link) => link.url === url) ?? false;
  const isSaved = adding.state === "added" || (!isSearching && hasSavedLink);
  const isAddingToEmptyList =
    adding.state === "adding" && links?.links.length === 0;
  const isLoading = isSearching || (!links && !error) || isAddingToEmptyList;

  function handleQuery(value: string) {
    void setFilters({ q: value || null, page: null });
  }

  function handleReset() {
    setSearch({ query: "", page: 1 });
    void setFilters({ q: null, page: null });
  }

  function handlePage(page: number) {
    setSearch({ query, page });
    void setFilters({ page }, { history: "push" });
  }

  async function handleSubmission(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSearch({ query, page: 1 });
    void setFilters({ page: null });

    if (!url) {
      return;
    }

    if (isSaved || requestRef.current?.url === url) {
      return;
    }

    const request = { url };
    requestRef.current = request;
    setAdding({ state: "adding" });

    let saved: LinkReference;
    try {
      saved = await linkCreate({ url });
    } catch (error) {
      if (requestRef.current !== request) {
        return;
      }

      requestRef.current = undefined;
      setAdding({ state: "error", error: deriveError(error) });
      return;
    }

    // Refresh every cached link search, including one opened while adding.
    await mutate((key) => Array.isArray(key) && key[0] === getLinkListKey()[0]);

    if (requestRef.current === request) {
      if (saved.url !== url) {
        setSearch({ query: saved.url, page: 1 });
        void setFilters({ q: saved.url, page: null });
        return;
      }

      setAdding({ state: "added" });
    }
  }

  return {
    data: {
      query: filters.q,
      page,
      isURL: url !== undefined,
      isSaved,
      isLoading,
      links,
      error,
      adding,
    },
    handlers: {
      handleQuery,
      handlePage,
      handleSubmission,
      handleReset,
      handleDismissError: () => setAdding({ state: "idle" }),
      handleRefresh: () => refresh(),
    },
  };
}
