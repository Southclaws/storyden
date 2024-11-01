"use client";

import { useCollectionList } from "@/api/openapi-client/collections";
import { Account, CollectionListOKResponse } from "@/api/openapi-schema";
import { CollectionCard } from "@/components/collection/CollectionCard";
import { CollectionCreateTrigger } from "@/components/content/CollectionCreate/CollectionCreateTrigger";
import { UnreadyBanner } from "@/components/site/Unready";
import { Breadcrumbs } from "@/components/ui/Breadcrumbs";
import { Heading } from "@/components/ui/heading";
import { CardGrid } from "@/components/ui/rich-card";
import { LStack } from "@/styled-system/jsx";

export type Props = {
  session?: Account;
  initialCollections: CollectionListOKResponse;
};

export function CollectionIndexScreen(props: Props) {
  const { data, error } = useCollectionList();
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
          <CollectionCreateTrigger
            session={props.session}
            size="xs"
            label="Create"
          />
        )}
      </Breadcrumbs>

      <CardGrid>
        {data.collections.map((collection) => (
          <CollectionCard key={collection.id} collection={collection} />
        ))}
      </CardGrid>
    </LStack>
  );
}
