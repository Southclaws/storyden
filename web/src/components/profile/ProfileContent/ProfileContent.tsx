import { Unready } from "src/components/site/Unready";

import { CollectionCard } from "@/components/collection/CollectionCard";
import { CollectionCreateTrigger } from "@/components/content/CollectionCreate/CollectionCreateTrigger";
import { ThreadReferenceList } from "@/components/post/ThreadReferenceList";
import { CardGrid } from "@/components/ui/rich-card";
import * as Tabs from "@/components/ui/tabs";
import { HStack, VStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { Props, useProfileContent } from "./useProfileContent";

export function ProfileContent(props: Props) {
  const { ready, error, data, isSelf } = useProfileContent(props);

  if (!ready) {
    return <Unready error={error} />;
  }

  const { threads, collections } = data;

  return (
    <VStack alignItems="start" w="full">
      <Tabs.Root width="full" variant="enclosed" defaultValue="threads">
        <Tabs.List>
          <Tabs.Trigger value="threads">Threads</Tabs.Trigger>
          <Tabs.Trigger value="collections">Collections</Tabs.Trigger>
          <Tabs.Indicator />
        </Tabs.List>

        <Tabs.Content value="threads">
          <ThreadReferenceList threads={threads} />
        </Tabs.Content>

        <Tabs.Content className={lstack()} value="collections">
          {isSelf && props.session && (
            <HStack w="full" justify="end">
              <CollectionCreateTrigger session={props.session} />
            </HStack>
          )}

          <CardGrid>
            {collections.map((collection) => (
              <CollectionCard
                key={collection.id}
                collection={collection}
                hideOwner
              />
            ))}
          </CardGrid>
        </Tabs.Content>
      </Tabs.Root>
    </VStack>
  );
}
