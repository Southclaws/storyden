import { dequal } from "dequal";
import { debounce } from "lodash/fp";
import { useQueryState } from "nuqs";
import { PropsWithChildren, createContext, useContext } from "react";

import { TagNameList, TagReference } from "@/api/openapi-schema";
import { SortState, useSortIndicator } from "@/components/site/SortIndicator";

type DirectoryBlockContextValue = {
  searchQuery: string;
  handleSearch: (q: string) => void;

  sort: SortState | null;
  handleSort: (property: string) => void;

  tagFilters: TagNameList;
  handleTagFilter: (tag: TagReference) => Promise<void>;
  highlightedTags: TagNameList | undefined;
};

export const LibraryPageDirectoryBlockContext = createContext<
  DirectoryBlockContextValue | undefined
>(undefined);

export function useDirectoryBlockContext() {
  const v = useContext(LibraryPageDirectoryBlockContext);
  if (!v) {
    throw new Error(
      "useDirectoryBlockContext must be used within a LibraryPageDirectoryBlockContextProvider",
    );
  }

  return v;
}

export function LibraryPageDirectoryBlockContextProvider({
  children,
}: PropsWithChildren) {
  const { sort, handleSort } = useSortIndicator();

  const [searchQuery, setSearchQuery] = useQueryState("search", {
    history: "replace",
    defaultValue: "",
    clearOnDefault: true,
  });

  const [tagFilters, setTagFilters] = useQueryState<string[]>("tag", {
    defaultValue: [],
    clearOnDefault: true,
    // This ensures the query params are removed entirely when tags are empty.
    eq: dequal,
    parse: (value) => {
      if (value === null || value === undefined) {
        return [];
      }
      return value.split(",").filter((v) => v.trim() !== "");
    },
  });

  async function handleTagFilter(tag: TagReference) {
    const present = tagFilters.includes(tag.name);
    if (present) {
      setTagFilters((prev) => prev.filter((t) => t !== tag.name));
    } else {
      setTagFilters((prev) => [...prev, tag.name]);
    }
  }

  // NOTE: We use `undefined` when the tag list is empty to cause the tag list
  // to render all tags as "highlighted" (not muted) so it's clearer.
  const highlightedTags = tagFilters.length > 0 ? tagFilters : undefined;

  const handleSearch = debounce(285, setSearchQuery as (q: string) => void);

  return (
    <LibraryPageDirectoryBlockContext.Provider
      value={{
        searchQuery,
        handleSearch,

        sort,
        handleSort,

        tagFilters,
        handleTagFilter,
        highlightedTags,
      }}
    >
      {children}
    </LibraryPageDirectoryBlockContext.Provider>
  );
}
