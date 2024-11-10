import { DatagraphSearchResult } from "src/api/openapi-schema";
import { EmptyState } from "src/components/site/EmptyState";

import { styled } from "@/styled-system/jsx";

import { DatagraphItemCard } from "../datagraph/DatagraphItemCard";

type Props = {
  result: DatagraphSearchResult;
};

export function DatagraphSearchResults({ result }: Props) {
  if (!result.items?.length) {
    return (
      <EmptyState>
        <p>No items were found.</p>
      </EmptyState>
    );
  }

  return (
    <styled.ol width="full" display="flex" flexDirection="column" gap="4">
      {result.items.map((v) => (
        <DatagraphItemCard key={v.ref.id} item={v} />
      ))}
    </styled.ol>
  );
}
