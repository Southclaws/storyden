"use client";

import { DatagraphSearchResults } from "src/components/search/DatagraphSearchResults";
import { UnreadyBanner } from "src/components/site/Unready";

import { useDatagraphSearch } from "@/api/openapi-client/datagraph";
import { DatagraphSearchOKResponse } from "@/api/openapi-schema";
import { Search } from "@/components/search/Search/Search";
import { useSearchQueryState } from "@/components/search/Search/useSearch";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { LStack } from "@/styled-system/jsx";

type Props = {
  query: string;
  page: number;
  initialResults: DatagraphSearchOKResponse;
};

export function SearchScreen(props: Props) {
  const [query, setQuery] = useSearchQueryState();

  const { data, error, isLoading } = useDatagraphSearch(
    {
      q: query,
      page: props.page.toString(),
    },
    {
      swr: {
        fallbackData: props.initialResults,
      },
    },
  );

  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  return (
    <LStack>
      <Search query={props.query} isLoading={isLoading} />

      <PaginationControls
        path="/search"
        params={{ q: query ?? "" }}
        currentPage={props.page}
        totalPages={data.total_pages}
        pageSize={data.page_size}
      />

      <DatagraphSearchResults result={data} />
    </LStack>
  );
}
