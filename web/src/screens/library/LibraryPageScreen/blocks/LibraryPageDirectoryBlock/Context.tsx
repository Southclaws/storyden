import { debounce } from "lodash/fp";
import { useQueryState } from "nuqs";
import { PropsWithChildren, createContext, useContext } from "react";
import { mutate } from "swr";

import { getNodeListChildrenKey } from "@/api/openapi-client/nodes";
import {
  NodeListOKResponse,
  TagNameList,
  TagReference,
} from "@/api/openapi-schema";
import { SortState, useSortIndicator } from "@/components/site/SortIndicator";
import { deepEqual } from "@/utils/equality";

import { useLibraryPageContext } from "../../Context";

type DirectoryBlockContextValue = {
  searchQuery: string;
  handleSearch: (q: string) => void;
  handleMutateChildren: (data?: NodeListOKResponse) => void;

  sort: SortState | null;
  handleSort: (property: string) => void;

  tagFilters: TagNameList;
  handleTagFilter: (tag: TagReference) => Promise<void>;
  highlightedTags: TagNameList | undefined;
  childrenSort: string | undefined;
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
  const { nodeID } = useLibraryPageContext();

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
    eq: deepEqual,
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

  // format the sort property as "name" or "-name" for asc/desc
  const childrenSort =
    sort !== null
      ? sort?.order === "asc"
        ? sort.property
        : `-${sort.property}`
      : undefined;

  const key = getNodeListChildrenKey(nodeID, {
    children_sort: childrenSort,
    tags: tagFilters,
    q: searchQuery,
  });

  function handleMutateChildren(data?: NodeListOKResponse) {
    mutate<NodeListOKResponse>(key, data);
  }

  return (
    <LibraryPageDirectoryBlockContext.Provider
      value={{
        searchQuery,
        handleSearch,
        handleMutateChildren,

        sort,
        handleSort,

        tagFilters,
        handleTagFilter,
        highlightedTags,
        childrenSort,
      }}
    >
      {children}
    </LibraryPageDirectoryBlockContext.Provider>
  );
}
