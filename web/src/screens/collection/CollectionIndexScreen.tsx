"use client";

import { parseAsInteger, useQueryStates } from "nuqs";

import { useCollectionList } from "@/api/openapi-client/collections";
import { Account, CollectionListOKResponse } from "@/api/openapi-schema";
import { CollectionCard } from "@/components/collection/CollectionCard";
import { CollectionCreateTrigger } from "@/components/content/CollectionCreate/CollectionCreateTrigger";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { UnreadyBanner } from "@/components/site/Unready";
import { Breadcrumbs } from "@/components/ui/breadcrumbs";
import { CardGrid } from "@/components/ui/rich-card";
import { LStack } from "@/styled-system/jsx";

export type Props = {
  session?: Account;
  initialCollections: CollectionListOKResponse;
};

export function CollectionIndexScreen(props: Props) {
  const [filters] = useQueryStates({
    page: parseAsInteger.withDefault(1),
  });

  const { data, error } = useCollectionList(
    { page: filters.page.toString() },
    { swr: { fallbackData: props.initialCollections } },
  );
  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  return (
    <LStack>
      <Breadcrumbs
        index={{
          href: "/c",
          label: "Collections",
        }}
        crumbs={[]}
      >
        {props.session && (
          <CollectionCreateTrigger session={props.session} label="Create" />
        )}
      </Breadcrumbs>

      <CardGrid>
        {data.collections.map((collection) => (
          <CollectionCard key={collection.id} collection={collection} />
        ))}
      </CardGrid>

      <PaginationControls
        path="/c"
        currentPage={filters.page}
        totalPages={data.total_pages}
        pageSize={data.page_size}
      />
    </LStack>
  );
}
