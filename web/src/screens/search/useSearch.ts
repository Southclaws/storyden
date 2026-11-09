"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { debounce } from "lodash";
import { parseAsArrayOf, parseAsString, useQueryStates } from "nuqs";
import { parseAsInteger, useQueryState } from "nuqs";
import { useEffect, useMemo, useRef, useState } from "react";
import { useForm } from "react-hook-form";

import { useCategoryList } from "@/api/openapi-client/categories";
import { useDatagraphSearch } from "@/api/openapi-client/datagraph";
import { useProfileList } from "@/api/openapi-client/profiles";
import { useTagList } from "@/api/openapi-client/tags";
import {
  DatagraphItemKind,
  DatagraphSearchOKResponse,
} from "@/api/openapi-schema";
import { useSearchQueryState } from "@/components/search/Search/useSearch";
import { MultiSelectPickerItem } from "@/components/ui/multi-select-picker";
import { DatagraphKindSchema } from "@/lib/datagraph/schema";
import * as z from "@/lib/validation";
import { deriveError } from "@/utils/error";

export type Props = {
  initialQuery: string;
  initialPage: number;
  initialKind: DatagraphItemKind[];
  initialAuthors: string[];
  initialCategories: string[];
  initialTags: string[];
  initialResults?: DatagraphSearchOKResponse;
};

export const FormSchema = z.object({
  q: z.string(),
  kind: z.optional(z.array(DatagraphKindSchema)),
});
export type Form = z.infer<typeof FormSchema>;

export function useSearchScreen(props: Props) {
  const [query, setQuery] = useSearchQueryState();
  const [page, setPage] = useQueryState("page", {
    ...parseAsInteger,
    defaultValue: 1,
  });

  const [filters, setFilters] = useQueryStates(
    {
      authors: parseAsArrayOf(parseAsString).withDefault(
        props.initialAuthors || [],
      ),
      categories: parseAsArrayOf(parseAsString).withDefault(
        props.initialCategories || [],
      ),
      tags: parseAsArrayOf(parseAsString).withDefault(props.initialTags || []),
      kind: parseAsArrayOf(parseAsString),
    },
    {
      history: "replace",
    },
  );

  const [authorsQuery, setAuthorsQuery] = useState("");
  const [categoriesQuery, setCategoriesQuery] = useState("");
  const [tagsQuery, setTagsQuery] = useState("");

  const {
    data: authorsData,
    error: authorsError,
    isLoading: authorsLoading,
  } = useProfileList({ q: authorsQuery });
  const {
    data: categoriesData,
    error: categoriesError,
    isLoading: categoriesLoading,
  } = useCategoryList();
  const {
    data: tagsData,
    error: tagsError,
    isLoading: tagsLoading,
  } = useTagList({ q: tagsQuery });

  const allCategories = useMemo(
    () => categoriesData?.categories ?? [],
    [categoriesData],
  );
  const authorsResults = (authorsData?.profiles ?? []).map((profile) => ({
    label: profile.handle,
    value: profile.handle,
  }));
  const lowerCategoriesQuery = categoriesQuery.toLowerCase();
  const categoriesResults = allCategories
    .filter(
      (category) =>
        category.name.toLowerCase().includes(lowerCategoriesQuery) ||
        category.slug.toLowerCase().includes(lowerCategoriesQuery),
    )
    .map((category) => ({ label: category.name, value: category.slug }));
  const tagsResults = (tagsData?.tags ?? []).map((tag) => ({
    label: tag.name,
    value: tag.name,
  }));

  const debouncedSearchTags = useRef(debounce(setTagsQuery, 300)).current;
  const debouncedSearchAuthors = useRef(debounce(setAuthorsQuery, 300)).current;

  useEffect(() => {
    return () => {
      debouncedSearchTags.cancel();
      debouncedSearchAuthors.cancel();
    };
  }, [debouncedSearchTags, debouncedSearchAuthors]);

  function handleQueryTags(query: string) {
    if (query.length === 0) {
      debouncedSearchTags.cancel();
      setTagsQuery("");
      return;
    }

    debouncedSearchTags(query);
  }

  function handleQueryCategories(query: string) {
    setCategoriesQuery(query);
  }

  function handleQueryAuthors(query: string) {
    if (query.length === 0) {
      debouncedSearchAuthors.cancel();
      setAuthorsQuery("");
      return;
    }

    debouncedSearchAuthors(query);
  }

  const handleAuthorsChange = async (items: MultiSelectPickerItem[]) => {
    const handles = items.map((item) => item.value);
    await setFilters({ authors: handles.length > 0 ? handles : null });
  };

  const handleCategoriesChange = async (items: MultiSelectPickerItem[]) => {
    const slugs = items.map((item) => item.value);
    await setFilters({ categories: slugs.length > 0 ? slugs : null });
  };

  const handleTagsChange = async (items: MultiSelectPickerItem[]) => {
    const tagNames = items.map((item) => item.value);
    await setFilters({ tags: tagNames.length > 0 ? tagNames : null });
  };

  const authorsValue = (filters.authors || []).map((v) => ({
    label: v,
    value: v,
  }));

  const categoriesValue = (filters.categories || []).map((slug) => {
    const category = allCategories.find((c) => c.slug === slug);
    return {
      label: category?.name || slug,
      value: slug,
    };
  });

  const tagsValue = (filters.tags || []).map((v) => ({
    label: v,
    value: v,
  }));

  const selectedKinds = filters.kind || [];
  const showCategories =
    selectedKinds.length === 0 || selectedKinds.includes("thread");
  const showTags =
    selectedKinds.length === 0 ||
    selectedKinds.includes("thread") ||
    selectedKinds.includes("node");

  const [categoryIds, setCategoryIds] = useState<string[]>([]);

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      q: props.initialQuery,
      kind: props.initialKind,
    },
  });

  const kind = form.watch("kind");

  useEffect(() => {
    if (!filters.categories || filters.categories.length === 0) {
      setCategoryIds([]);
      return;
    }

    const ids = filters.categories
      .map((slug) => {
        const cat = allCategories.find((c) => c.slug === slug);
        return cat?.id;
      })
      .filter((id): id is string => !!id);

    setCategoryIds(ids);
  }, [filters.categories, allCategories]);

  useEffect(() => {
    const url = new URL(window.location.href);

    url.searchParams.delete("kind");
    if (kind && kind.length > 0) {
      kind.forEach((k) => url.searchParams.append("kind", k));
    }

    window.history.replaceState({}, "", url.toString());
  }, [kind]);

  const hasActiveFilters =
    (filters.authors?.length ?? 0) > 0 ||
    (filters.categories?.length ?? 0) > 0 ||
    (filters.tags?.length ?? 0) > 0;

  const { data, error, isLoading, mutate } = useDatagraphSearch(
    {
      q: query,
      kind: kind,
      page: page.toString(),
      authors: filters.authors || undefined,
      categories: categoryIds.length > 0 ? categoryIds : undefined,
      tags: filters.tags || undefined,
    },
    {
      swr: {
        enabled: !!query || hasActiveFilters,
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
    await mutate({
      current_page: 1,
      total_pages: 0,
      results: 0,
      items: [],
      page_size: 50,
    });
    form.clearErrors();
    form.reset();
    setQuery("");
  };

  return {
    ready: true as const,
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
      handleQueryAuthors,
      handleQueryCategories,
      handleQueryTags,
      handleAuthorsChange,
      handleCategoriesChange,
      handleTagsChange,
    },
    filters: {
      authorsValue,
      authorsResults,
      authorsError: deriveError(authorsError),
      authorsLoading,
      categoriesValue,
      categoriesResults,
      categoriesError: deriveError(categoriesError),
      categoriesLoading,
      tagsValue,
      tagsResults,
      tagsError: deriveError(tagsError),
      tagsLoading,
      showCategories,
      showTags,
    },
  };
}
