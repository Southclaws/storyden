"use client";

import { debounce } from "lodash";
import { usePathname } from "next/navigation";
import { ChangeEvent, useEffect, useState } from "react";

import { postSearch } from "src/api/openapi/posts";
import { PostProps } from "src/api/openapi/schemas";
import { useSession } from "src/auth";

export function useNavpill() {
  const pathname = usePathname();
  const [isExpanded, setExpanded] = useState(false);
  const account = useSession();
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<PostProps[]>([]);

  // Close the menu for either navigation events or outside clicks/taps:

  useEffect(() => setExpanded(false), [pathname]);

  function onExpand() {
    setExpanded(true);
  }

  function onClose() {
    setExpanded(false);
  }

  const doSearch = debounce(async (v: string) => {
    postSearch({ body: v })
      .then((results) => setSearchResults(results.results))
      .catch((e) => {
        console.log({ e });
      });
  }, 250);

  async function onSearch(e: ChangeEvent<HTMLInputElement>) {
    const query = e.target.value;

    setSearchQuery(query);

    if (query === "") {
      setSearchResults([]);
      return;
    }

    await doSearch(e.target.value);
  }

  return {
    isExpanded,
    onExpand,
    onClose,
    account,
    searchQuery,
    onSearch,
    searchResults,
  };
}
