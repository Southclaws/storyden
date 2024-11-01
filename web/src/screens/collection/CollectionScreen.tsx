"use client";

import { CollectionWithItems } from "src/api/openapi-schema";
import { Unready } from "src/components/site/Unready";

import { useCollectionGet } from "@/api/openapi-client/collections";
import { Account } from "@/api/openapi-schema";
import { CollectionCreateTrigger } from "@/components/content/CollectionCreate/CollectionCreateTrigger";
import { DatagraphItemCard } from "@/components/datagraph/DatagraphItemCard";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Breadcrumbs } from "@/components/ui/Breadcrumbs";
import { CardGrid } from "@/components/ui/rich-card";
import { LStack, VStack, styled } from "@/styled-system/jsx";

type Props = {
  session?: Account;
  initialCollection: CollectionWithItems;
};

export function CollectionScreen({ session, initialCollection }: Props) {
  const { data, error } = useCollectionGet(initialCollection.id, {
    swr: { fallbackData: initialCollection },
  });
  if (!data) {
    return <Unready error={error} />;
  }

  const collection = data;

  const url = `/c/${collection.id}`;

  return (
    <VStack alignItems="start">
      <Breadcrumbs
        index={{
          href: "/c",
          label: "Collections",
        }}
        crumbs={[{ label: collection.name, href: url }]}
      >
        {session && (
          <CollectionCreateTrigger session={session} size="xs" label="Create" />
        )}
      </Breadcrumbs>

      <LStack>
        <styled.p fontSize="sm">{collection.description}</styled.p>
        <MemberBadge
          profile={collection.owner}
          name="full-horizontal"
          size="sm"
        />
      </LStack>

      <CardGrid>
        {collection.items.map((i) => (
          <DatagraphItemCard key={i.id} item={i} />
        ))}
      </CardGrid>
    </VStack>
  );
}
