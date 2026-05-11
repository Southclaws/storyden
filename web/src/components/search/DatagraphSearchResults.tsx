import { DatagraphSearchResult } from "src/api/openapi-schema";
import { EmptyState } from "src/components/site/EmptyState";

import { useI18n } from "@/i18n/provider";
import { styled } from "@/styled-system/jsx";

import { DatagraphItemCard } from "../datagraph/DatagraphItemCard";

type Props = {
  result: DatagraphSearchResult;
};

export function DatagraphSearchResults({ result }: Props) {
  const { t } = useI18n();

  if (!result.items?.length) {
    return (
      <EmptyState>
        <p>{t("No items were found.")}</p>
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
