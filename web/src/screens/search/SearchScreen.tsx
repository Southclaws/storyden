"use client";

import { DatagraphSearchResults } from "src/components/search/DatagraphSearchResults";
import { UnreadyBanner } from "src/components/site/Unready";

import { useDatagraphSearch } from "@/api/openapi-client/datagraph";
import {
  DatagraphItemKind,
  DatagraphSearchOKResponse,
} from "@/api/openapi-schema";
import { Search } from "@/components/search/Search/Search";
import { useSearchQueryState } from "@/components/search/Search/useSearch";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { VStack } from "@/styled-system/jsx";

type Props = {
  query: string;
  page: number;
  kind: DatagraphItemKind[];
  initialResults: DatagraphSearchOKResponse;
};

export function SearchScreen(props: Props) {
  const [query, setQuery] = useSearchQueryState();

  const { data, error, isLoading } = useDatagraphSearch(
    {
      q: query,
      page: props.page.toString(),
      kind: props.kind,
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
    <VStack>
      <Search query={props.query} isLoading={isLoading} />

      <PaginationControls
        path="/search"
        params={{ q: query ?? "" }}
        currentPage={props.page}
        totalPages={data.total_pages}
        pageSize={data.page_size}
      />

      <DatagraphSearchResults result={data} />
    </VStack>
  );
}
