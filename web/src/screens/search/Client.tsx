"use client";

import { DatagraphSearchResult } from "src/api/openapi/schemas";
import { DatagraphSearchResults } from "src/components/search/DatagraphSearchResults";
import { useSearch } from "src/components/search/useSearch";
import { Unready } from "src/components/site/Unready";

export function Client(props: {
  query: string;
  results: DatagraphSearchResult;
}) {
  const { data, error } = useSearch(
    {
      q: props.query,
    },
    props.results,
  );

  if (!data) return <Unready {...error} />;

  return <DatagraphSearchResults result={data} />;
}
