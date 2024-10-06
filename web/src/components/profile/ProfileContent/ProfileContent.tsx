import { Unready } from "src/components/site/Unready";

import { ThreadItemList } from "@/components/feed/ThreadItemList";
import * as Tabs from "@/components/ui/tabs";
import { VStack } from "@/styled-system/jsx";

import { CollectionList } from "../CollectionList/CollectionList";

import { Props, useProfileContent } from "./useProfileContent";

export function ProfileContent(props: Props) {
  const content = useProfileContent(props);

  if (!content.ready) {
    return <Unready {...content.error} />;
  }

  return (
    <VStack alignItems="start" w="full">
      <Tabs.Root width="full" variant="line" defaultValue="threads">
        <Tabs.List>
          <Tabs.Trigger value="threads">Threads</Tabs.Trigger>
          <Tabs.Trigger value="collections">Collections</Tabs.Trigger>
          <Tabs.Indicator />
        </Tabs.List>

        <Tabs.Content value="threads">
          <ThreadItemList threads={content.data.threads} />
        </Tabs.Content>

        <Tabs.Content value="collections">
          <CollectionList collections={content.data.collections} />
        </Tabs.Content>
      </Tabs.Root>
    </VStack>
  );
}
