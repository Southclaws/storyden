import { debounce } from "lodash/fp";
import { useQueryState } from "nuqs";
import { PropsWithChildren, createContext, useContext } from "react";

export const LibraryPageDirectoryBlockContext = createContext({
  searchQuery: "",
  handleSearch: (_: string) => {},
});

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
  const [searchQuery, setSearchQuery] = useQueryState("search", {
    history: "replace",
    defaultValue: "",
    clearOnDefault: true,
  });

  const handleSearch = debounce(285, setSearchQuery as (q: string) => void);

  return (
    <LibraryPageDirectoryBlockContext.Provider
      value={{
        searchQuery,
        handleSearch,
      }}
    >
      {children}
    </LibraryPageDirectoryBlockContext.Provider>
  );
}
