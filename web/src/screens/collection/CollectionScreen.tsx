"use client";

import { useCollectionGet } from "@/api/openapi-client/collections";
import { CollectionWithItems } from "@/api/openapi-schema";
import { Account } from "@/api/openapi-schema";
import { CollectionCreateTrigger } from "@/components/content/CollectionCreate/CollectionCreateTrigger";
import { DatagraphItemCard } from "@/components/datagraph/DatagraphItemCard";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Unready } from "@/components/site/Unready";
import { Breadcrumbs } from "@/components/ui/breadcrumbs";
import { PageHeading } from "@/components/ui/page-heading";
import { CardGrid } from "@/components/ui/rich-card";
import { Text } from "@/components/ui/text";
import { LStack, VStack } from "@/styled-system/jsx";

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

  const url = `/c/${collection.slug}`;

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
          <CollectionCreateTrigger session={session} label="Create" />
        )}
      </Breadcrumbs>

      <LStack gap="1">
        <PageHeading>{collection.name}</PageHeading>

        <Text variant="supporting" color="text.default">
          {collection.description ? (
            <span>{collection.description}</span>
          ) : (
            <Text as="span" variant="metadata" fontStyle="italic">
              (no description)
            </Text>
          )}
        </Text>

        <MemberBadge
          profile={collection.owner}
          name="full-horizontal"
          size="sm"
        />
      </LStack>

      <CardGrid>
        {collection.items.map((i) => (
          <DatagraphItemCard key={i.id} item={i.item} />
        ))}
      </CardGrid>
    </VStack>
  );
}
