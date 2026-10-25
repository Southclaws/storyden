import { DatagraphSearchResult } from "@/api/openapi-schema";
import { EmptyState } from "@/components/site/EmptyState";
import { styled } from "@/styled-system/jsx";

import { DatagraphItemCard } from "../datagraph/DatagraphItemCard";

import { useSearchHighlights } from "./useSearchHighlights";

type Props = {
  result: DatagraphSearchResult;
  query: string;
};

export function DatagraphSearchResults({ result, query }: Props) {
  const highlightRef = useSearchHighlights(query);
  if (!result.items?.length) {
    return (
      <EmptyState>
        <p>No items were found.</p>
      </EmptyState>
    );
  }

  return (
    <styled.ol
      ref={highlightRef}
      width="full"
      display="flex"
      flexDirection="column"
      gap="4"
    >
      {result.items.map((v) => (
        <DatagraphItemCard key={v.ref.id} item={v} />
      ))}
    </styled.ol>
  );
}
