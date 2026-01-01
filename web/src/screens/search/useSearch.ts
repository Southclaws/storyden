"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { debounce } from "lodash";
import { parseAsArrayOf, parseAsString, useQueryStates } from "nuqs";
import { parseAsInteger, useQueryState } from "nuqs";
import { useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { categoryList } from "@/api/openapi-client/categories";
import { useDatagraphSearch } from "@/api/openapi-client/datagraph";
import { profileList } from "@/api/openapi-client/profiles";
import { tagList } from "@/api/openapi-client/tags";
import { Category } from "@/api/openapi-schema";
import {
  DatagraphItemKind,
  DatagraphSearchOKResponse,
} from "@/api/openapi-schema";
import { useSearchQueryState } from "@/components/search/Search/useSearch";
import { MultiSelectPickerItem } from "@/components/ui/MultiSelectPicker";
import { DatagraphKindSchema } from "@/lib/datagraph/schema";
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

  const [allCategories, setAllCategories] = useState<Category[]>([]);
  const [authorsResults, setAuthorsResults] = useState<MultiSelectPickerItem[]>(
    [],
  );
  const [categoriesResults, setCategoriesResults] = useState<
    MultiSelectPickerItem[]
  >([]);
  const [tagsResults, setTagsResults] = useState<MultiSelectPickerItem[]>([]);

  const [authorsError, setAuthorsError] = useState<string | null>(null);
  const [categoriesError, setCategoriesError] = useState<string | null>(null);
  const [tagsError, setTagsError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchCategories() {
      try {
        const { categories } = await categoryList();
        setAllCategories(categories);
      } catch (error) {
        setCategoriesError(deriveError(error));
      }
    }

    fetchCategories();
  }, []);

  const debouncedSearchTags = useRef(
    debounce(async (query: string) => {
      try {
        setTagsError(null);

        if (query.length === 0) {
          setTagsResults([]);
          return;
        }

        const { tags } = await tagList({ q: query });
        const items = tags.map((t) => ({ label: t.name, value: t.name }));
        setTagsResults(items);
      } catch (error) {
        setTagsError(deriveError(error));
        setTagsResults([]);
      }
    }, 300),
  ).current;

  const debouncedSearchAuthors = useRef(
    debounce(async (query: string) => {
      try {
        setAuthorsError(null);

        if (query.length === 0) {
          setAuthorsResults([]);
          return;
        }

        const { profiles } = await profileList({ q: query });
        const items = profiles.map((v) => ({
          label: v.handle,
          value: v.handle,
        }));
        setAuthorsResults(items);
      } catch (error) {
        setAuthorsError(deriveError(error));
        setAuthorsResults([]);
      }
    }, 300),
  ).current;

  useEffect(() => {
    return () => {
      debouncedSearchTags.cancel();
      debouncedSearchAuthors.cancel();
    };
  }, []);

  function handleQueryTags(query: string) {
    if (query.length === 0) {
      setTagsResults([]);
    } else {
      debouncedSearchTags(query);
    }
  }

  function handleQueryCategories(query: string) {
    if (query.length === 0) {
      setCategoriesResults([]);
      return;
    }

    const filtered = allCategories
      .filter(
        (cat) =>
          cat.name.toLowerCase().includes(query.toLowerCase()) ||
          cat.slug.toLowerCase().includes(query.toLowerCase()),
      )
      .map((cat) => ({ label: cat.name, value: cat.slug }));

    setCategoriesResults(filtered);
  }

  function handleQueryAuthors(query: string) {
    if (query.length === 0) {
      setAuthorsResults([]);
    } else {
      debouncedSearchAuthors(query);
    }
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
      authorsError,
      categoriesValue,
      categoriesResults,
      categoriesError,
      tagsValue,
      tagsResults,
      tagsError,
      showCategories,
      showTags,
    },
  };
}
